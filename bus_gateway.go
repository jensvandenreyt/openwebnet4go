package openwebnet

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jensvandenreyt/openwebnet4go/communication"
	"github.com/jensvandenreyt/openwebnet4go/message"
	"github.com/rs/zerolog/log"
)

const (
	defaultBUSPort      = 20000
	reconnectRetryAfter = 2500 * time.Millisecond
	reconnectRetryMax   = 60 * time.Second
	reconnectMultiplier = 2
	connectionTimeoutMS = 120000
)

// BUSGateway connects to BUS OpenWebNet gateways.
type BUSGateway struct {
	host string
	port int
	pwd  string

	connector       *communication.BUSConnector
	isConnected     bool
	isDiscovering   bool
	closedRequested bool

	macAddr         []byte
	firmwareVersion string

	listeners []GatewayListener
	mu        sync.Mutex
}

// NewBUSGateway creates a new BUSGateway with the given host, port and password.
func NewBUSGateway(host string, port int, pwd string) Gateway {
	return &BUSGateway{
		host:      host,
		port:      port,
		pwd:       pwd,
		listeners: make([]GatewayListener, 0),
	}
}

// GetHost returns the gateway host.
func (gw *BUSGateway) GetHost() string { return gw.host }

// GetPort returns the gateway port.
func (gw *BUSGateway) GetPort() int { return gw.port }

// GetPassword returns the gateway password.
func (gw *BUSGateway) GetPassword() string { return gw.pwd }

// IsConnected returns true if the gateway is connected.
func (gw *BUSGateway) IsConnected() bool { return gw.isConnected }

// Connect connects to the OpenWebNet gateway.
func (gw *BUSGateway) Connect(ctx context.Context) error {
	if gw.isConnected {
		log.Trace().Msg("##BUS## OpenGateway is already connected")
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	gw.closedRequested = false
	gw.connector = communication.NewBUSConnector(gw.host, gw.port, gw.pwd)
	gw.connector.SetListener(gw)

	if err := gw.connector.OpenMonConn(); err != nil {
		gw.notifyListeners(func(l GatewayListener) { l.OnConnectionError(err) })
		return err
	}
	if !gw.connector.IsMonConnected() {
		err := fmt.Errorf("MON connection failed")
		gw.notifyListeners(func(l GatewayListener) { l.OnConnectionError(err) })
		return err
	}

	if err := gw.connector.OpenCmdConn(); err != nil {
		gw.notifyListeners(func(l GatewayListener) { l.OnConnectionError(err) })
		return err
	}
	if !gw.connector.IsCmdConnected() {
		err := fmt.Errorf("CMD connection failed")
		gw.notifyListeners(func(l GatewayListener) { l.OnConnectionError(err) })
		return err
	}

	// Request MAC address and firmware version
	gw.handleManagementDimensions(gw.sendInternal(message.GatewayMgmtRequestMACAddress()))
	gw.handleManagementDimensions(gw.sendInternal(message.GatewayMgmtRequestFirmwareVersion()))

	log.Trace().Msg("##GW## ============ OpenGateway CONNECTED! ============")
	gw.isConnected = true

	// Tie the connection lifecycle to the provided context.
	go func() {
		<-ctx.Done()
		gw.CloseConnection()
	}()

	gw.notifyListeners(func(l GatewayListener) { l.OnConnected() })
	return nil
}

// Reconnect tries to reconnect to the gateway with increasing retry intervals.
func (gw *BUSGateway) Reconnect() error {
	retry := reconnectRetryAfter
	for !gw.isConnected && !gw.closedRequested {
		log.Trace().Msgf("--Sleeping %v before re-connecting...", retry)
		time.Sleep(retry)
		if gw.closedRequested {
			break
		}

		gw.connector = communication.NewBUSConnector(gw.host, gw.port, gw.pwd)
		gw.connector.SetListener(gw)

		err := gw.connector.OpenMonConn()
		if err != nil {
			if _, ok := err.(*communication.OWNAuthError); ok {
				return err
			}
			retry = retry * time.Duration(reconnectMultiplier)
			if retry > reconnectRetryMax {
				retry = reconnectRetryMax
			}
			gw.notifyListeners(func(l GatewayListener) { l.OnConnectionError(err) })
			continue
		}

		if !gw.connector.IsMonConnected() {
			continue
		}

		err = gw.connector.OpenCmdConn()
		if err != nil {
			if _, ok := err.(*communication.OWNAuthError); ok {
				return err
			}
			retry = retry * time.Duration(reconnectMultiplier)
			if retry > reconnectRetryMax {
				retry = reconnectRetryMax
			}
			gw.notifyListeners(func(l GatewayListener) { l.OnConnectionError(err) })
			continue
		}

		if gw.connector.IsCmdConnected() {
			gw.handleManagementDimensions(gw.sendInternal(message.GatewayMgmtRequestMACAddress()))
			gw.handleManagementDimensions(gw.sendInternal(message.GatewayMgmtRequestFirmwareVersion()))
			gw.isConnected = true
			gw.notifyListeners(func(l GatewayListener) { l.OnReconnected() })
		}
	}
	return nil
}

// Send sends a command message and returns the response messages.
func (gw *BUSGateway) Send(msg message.OpenMessage) (*communication.Response, error) {
	if !gw.isConnected {
		return nil, communication.NewOWNError("Error while sending message: the gateway is not connected")
	}
	return gw.sendInternal(msg)
}

func (gw *BUSGateway) sendInternal(msg message.OpenMessage) (*communication.Response, error) {
	return gw.connector.SendCommandSynch(msg.GetFrameValue())
}

// Subscribe adds a listener for gateway events.
func (gw *BUSGateway) Subscribe(listener GatewayListener) {
	gw.mu.Lock()
	defer gw.mu.Unlock()
	for _, l := range gw.listeners {
		if l == listener {
			return
		}
	}
	gw.listeners = append(gw.listeners, listener)
}

// Unsubscribe removes a listener.
func (gw *BUSGateway) Unsubscribe(listener GatewayListener) {
	gw.mu.Lock()
	defer gw.mu.Unlock()
	for i, l := range gw.listeners {
		if l == listener {
			gw.listeners = append(gw.listeners[:i], gw.listeners[i+1:]...)
			return
		}
	}
}

func (gw *BUSGateway) notifyListeners(method func(GatewayListener)) {
	gw.mu.Lock()
	listenersCopy := make([]GatewayListener, len(gw.listeners))
	copy(listenersCopy, gw.listeners)
	gw.mu.Unlock()

	go func() {
		for _, l := range listenersCopy {
			method(l)
		}
	}()
}

func (gw *BUSGateway) handleManagementDimensions(res *communication.Response, err error) {
	if err != nil || res == nil {
		return
	}
	for _, msg := range res.GetResponseMessages() {
		bom, ok := msg.(*message.BaseOpenMessage)
		if !ok || !message.IsGatewayMgmtMessage(bom) {
			continue
		}
		dim := bom.GetDim()
		if dim == nil {
			continue
		}
		switch dim.Value() {
		case message.DimGatewayMgmtMACAddress.Value():
			mac, err := message.GatewayMgmtParseMACAddress(bom)
			if err == nil {
				gw.macAddr = mac
				log.Trace().Msgf("##GW## MAC ADDRESS: %s", gw.GetMACAddr())
			}
		case message.DimGatewayMgmtFirmwareVersion.Value():
			fw, err := message.GatewayMgmtParseFirmwareVersion(bom)
			if err == nil {
				gw.firmwareVersion = fw
				log.Trace().Msgf("##GW## FIRMWARE: %s", gw.GetFirmwareVersion())
			}
		}
	}
}

// GetFirmwareVersion returns the firmware version.
func (gw *BUSGateway) GetFirmwareVersion() string {
	return gw.firmwareVersion
}

// GetMACAddr returns the MAC address as human-readable string.
func (gw *BUSGateway) GetMACAddr() string {
	return message.GatewayMgmtFormatMACAddress(gw.macAddr)
}

// CloseConnection closes the connection to the gateway.
func (gw *BUSGateway) CloseConnection() {
	gw.closedRequested = true
	if gw.connector != nil {
		gw.connector.Disconnect()
	}
	gw.isConnected = false
}

// IsCmdConnectionReady returns true if CMD connection is ready.
func (gw *BUSGateway) IsCmdConnectionReady() bool {
	if gw.isConnected && gw.connector.IsCmdConnected() {
		now := time.Now().UnixMilli()
		if now-gw.connector.GetLastCmdFrameSentTs() < connectionTimeoutMS {
			return true
		}
	}
	return false
}

// DiscoverDevices starts a device discovery session.
func (gw *BUSGateway) DiscoverDevices() error {
	if gw.isDiscovering {
		log.Trace().Msg("##BUS## discovery already in progress -> SKIPPING...")
		return nil
	}
	if !gw.isConnected {
		return fmt.Errorf("cannot perform discovery: gateway is not connected")
	}
	gw.isDiscovering = true
	return gw.discoverDevicesInternal()
}

func (gw *BUSGateway) discoverDevicesInternal() error {
	defer func() {
		gw.isDiscovering = false
	}()

	// DISCOVER LIGHTS
	log.Trace().Msg("##BUS## ----- LIGHTS discovery -----")
	res, err := gw.sendInternal(message.LightingRequestStatus(message.WhereLightAutomGeneral.Value()))
	if err == nil && res != nil {
		for _, msg := range res.GetResponseMessages() {
			if bom, ok := msg.(*message.BaseOpenMessage); ok && bom.WhoField == message.WhoLighting {
				if bom.DetectDeviceTyp != nil {
					devType, _ := bom.DetectDeviceTyp(bom)
					if devType != 0 {
						w := bom.GetWhere()
						gw.notifyListeners(func(l GatewayListener) { l.OnNewDevice(w, devType, bom) })
					}
				}
			}
		}
	}

	// DISCOVER AUTOMATION
	log.Trace().Msg("##BUS## ----- AUTOMATION discovery -----")
	res, err = gw.sendInternal(message.AutomationRequestStatus(message.WhereLightAutomGeneral.Value()))
	if err == nil && res != nil {
		for _, msg := range res.GetResponseMessages() {
			if bom, ok := msg.(*message.BaseOpenMessage); ok && bom.WhoField == message.WhoAutomation {
				if bom.DetectDeviceTyp != nil {
					devType, _ := bom.DetectDeviceTyp(bom)
					if devType != 0 {
						w := bom.GetWhere()
						gw.notifyListeners(func(l GatewayListener) { l.OnNewDevice(w, devType, bom) })
					}
				}
			}
		}
	}

	// DISCOVER ENERGY MANAGEMENT
	log.Trace().Msg("##BUS## ----- ENERGY MANAGEMENT discovery -----")
	res, err = gw.sendInternal(message.EnergyMgmtDiagnosticRequestDiagnostic(message.WhereEnergyManagementGeneral.Value()))
	if err == nil && res != nil {
		for _, msg := range res.GetResponseMessages() {
			if bom, ok := msg.(*message.BaseOpenMessage); ok && bom.WhoField == message.WhoEnergyManagementDiagnostic {
				if bom.DetectDeviceTyp != nil {
					devType, _ := bom.DetectDeviceTyp(bom)
					if devType != 0 {
						w := bom.GetWhere()
						gw.notifyListeners(func(l GatewayListener) { l.OnNewDevice(w, devType, bom) })
					}
				}
			}
		}
	}

	// DISCOVER THERMOREGULATION
	log.Trace().Msg("##BUS## ----- THERMOREGULATION discovery -----")
	res, err = gw.sendInternal(message.ThermoregulationDiagnosticRequestDiagnostic(message.WhereThermoAllMasterProbes.Value()))
	if err == nil && res != nil {
		for _, msg := range res.GetResponseMessages() {
			if bom, ok := msg.(*message.BaseOpenMessage); ok && bom.WhoField == message.WhoThermoregulationDiagnostic {
				if bom.DetectDeviceTyp != nil {
					devType, _ := bom.DetectDeviceTyp(bom)
					if devType != 0 {
						w := bom.GetWhere()
						gw.notifyListeners(func(l GatewayListener) { l.OnNewDevice(w, devType, bom) })
					}
				}
			}
		}
	}

	// DISCOVER DRY CONTACT / IR SENSOR
	log.Trace().Msg("##BUS## ----- DRY CONTACT / IR sensor discovery -----")
	res, err = gw.sendInternal(message.CENPlusRequestStatus("30"))
	if err == nil && res != nil {
		for _, msg := range res.GetResponseMessages() {
			if bom, ok := msg.(*message.BaseOpenMessage); ok && bom.WhoField == message.WhoCENPlusScenarioScheduler {
				if bom.DetectDeviceTyp != nil {
					devType, _ := bom.DetectDeviceTyp(bom)
					if devType != 0 {
						w := bom.GetWhere()
						gw.notifyListeners(func(l GatewayListener) { l.OnNewDevice(w, devType, bom) })
					}
				}
			}
		}
	}

	log.Trace().Msg("##BUS## ----- ### DISCOVERY COMPLETED")
	gw.notifyListeners(func(l GatewayListener) { l.OnDiscoveryCompleted() })
	return nil
}

// OnMessage implements ConnectorListener - called when a message is received on MON.
func (gw *BUSGateway) OnMessage(msg message.OpenMessage) {
	gw.notifyListeners(func(l GatewayListener) { l.OnEventMessage(msg) })
}

// OnMonDisconnected implements ConnectorListener - called when MON is disconnected.
func (gw *BUSGateway) OnMonDisconnected(err error) {
	gw.isConnected = false
	gw.notifyListeners(func(l GatewayListener) { l.OnDisconnected(err) })
}

// String returns a string representation.
func (gw *BUSGateway) String() string {
	return fmt.Sprintf("BUS_%s:%d", gw.host, gw.port)
}

// FormatMACAddress is a helper to format MAC bytes as XX:XX:XX... string.
func FormatMACAddress(mac []byte) string {
	if mac == nil {
		return ""
	}
	parts := make([]string, len(mac))
	for i, b := range mac {
		parts[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(parts, ":")
}
