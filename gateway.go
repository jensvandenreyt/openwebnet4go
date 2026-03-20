package openwebnet

import (
	"github.com/jensvandenreyt/openwebnet4go/communication"
	message2 "github.com/jensvandenreyt/openwebnet4go/message"
)

// GatewayListener defines callbacks for gateway events.
type GatewayListener interface {
	// OnConnected is called after the connection has been established.
	OnConnected()
	// OnConnectionError is called when connecting returns an error.
	OnConnectionError(err error)
	// OnConnectionClosed is called after the gateway connection has been closed.
	OnConnectionClosed()
	// OnDisconnected is called after the connection has been lost.
	OnDisconnected(err error)
	// OnReconnected is called after the connection has been re-connected.
	OnReconnected()
	// OnEventMessage is called when a new OpenWebNet message is received on the MON session.
	OnEventMessage(msg message2.OpenMessage)
	// OnNewDevice is called when a new device is discovered.
	OnNewDevice(where message2.Where, deviceType message2.OpenDeviceType, msg *message2.BaseOpenMessage)
	// OnDiscoveryCompleted is called when device discovery has been completed.
	OnDiscoveryCompleted()
}

// Gateway is the interface for OpenWebNet gateways.
type Gateway interface {
	// Connect connects to the gateway.
	Connect() error
	// Send sends a command message and returns the response.
	Send(msg message2.OpenMessage) (*communication.Response, error)
	// Subscribe adds a listener for gateway events.
	Subscribe(listener GatewayListener)
	// Unsubscribe removes a listener.
	Unsubscribe(listener GatewayListener)
	// DiscoverDevices starts a device discovery session.
	DiscoverDevices() error
	// IsConnected returns true if connected.
	IsConnected() bool
	// GetFirmwareVersion returns the gateway firmware version.
	GetFirmwareVersion() string
	// GetMACAddr returns the gateway MAC address as a human-readable string.
	GetMACAddr() string
	// CloseConnection closes the gateway connection.
	CloseConnection()
}
