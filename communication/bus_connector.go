package communication

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"time"

	message2 "github.com/jensvandenreyt/openwebnet4go/message"
)

// BUS connection constants.
const (
	MonType       = "MON"
	CmdType       = "CMD"
	CmdSession    = "*99*0##"
	CmdSessionAlt = "*99*9##"
	MonSession    = "*99*1##"

	SocketConnectTimeout = 5 * time.Second
	MonSocketReadTimeout = 120 * time.Second
	CmdSocketReadTimeout = 30 * time.Second
	MonKeepaliveTimer    = 90 * time.Second
	HandshakeTimeout     = 2 * time.Second

	HMACSHA1 = "*98*1##"
	HMACSHA2 = "*98*2##"
)

// ConnectorListener defines methods to receive MONITOR messages.
type ConnectorListener interface {
	OnMessage(msg message2.OpenMessage)
	OnMonDisconnected(err error)
}

// BUSConnector communicates with a BUS OpenWebNet gateway.
type BUSConnector struct {
	host string
	port int
	pwd  string

	cmdChannel *FrameChannel
	monChannel *FrameChannel
	cmdConn    net.Conn
	monConn    net.Conn

	isCmdConnected     bool
	isMonConnected     bool
	lastCmdFrameSentTs int64

	listener      ConnectorListener
	monStopChan   chan struct{}
	keepaliveDone chan struct{}
}

// NewBUSConnector creates a new BUSConnector.
func NewBUSConnector(host string, port int, pwd string) *BUSConnector {
	return &BUSConnector{
		host: host,
		port: port,
		pwd:  pwd,
	}
}

// SetListener sets the ConnectorListener for MONITOR events.
func (bc *BUSConnector) SetListener(listener ConnectorListener) {
	bc.listener = listener
}

// IsCmdConnected returns true if CMD is connected.
func (bc *BUSConnector) IsCmdConnected() bool {
	return bc.isCmdConnected
}

// IsMonConnected returns true if MON is connected.
func (bc *BUSConnector) IsMonConnected() bool {
	return bc.isMonConnected
}

// GetLastCmdFrameSentTs returns timestamp of last CMD frame sent.
func (bc *BUSConnector) GetLastCmdFrameSentTs() int64 {
	return bc.lastCmdFrameSentTs
}

// OpenCmdConn opens the command (CMD) connection.
func (bc *BUSConnector) OpenCmdConn() error {
	if bc.isCmdConnected {
		return nil
	}
	if err := bc.openConnection(CmdType); err != nil {
		return err
	}
	bc.isCmdConnected = true
	log.Println("##BUS-conn## ============ CMD CONNECTED ============")
	return nil
}

// OpenMonConn opens the monitor (MON) connection.
func (bc *BUSConnector) OpenMonConn() error {
	if bc.isMonConnected {
		return nil
	}
	if err := bc.openConnection(MonType); err != nil {
		return err
	}
	bc.isMonConnected = true
	log.Println("##BUS-conn## ============ MON CONNECTED ============")

	bc.monStopChan = make(chan struct{})
	go bc.monReceiveLoop()
	bc.startMonKeepalive()

	return nil
}

func (bc *BUSConnector) openConnection(connType string) error {
	log.Printf("##BUS-conn## Establishing %s connection to BUS Gateway on %s:%d...", connType, bc.host, bc.port)
	ch, err := bc.connectSocket(connType)
	if err != nil {
		return NewOWNErrorWithCause(
			fmt.Sprintf("Could not open BUS-%s connection to %s:%d", connType, bc.host, bc.port), err)
	}
	if err := bc.doHandshake(ch, connType); err != nil {
		return err
	}
	return nil
}

func (bc *BUSConnector) connectSocket(connType string) (*FrameChannel, error) {
	addr := net.JoinHostPort(bc.host, fmt.Sprintf("%d", bc.port))
	conn, err := net.DialTimeout("tcp", addr, SocketConnectTimeout)
	if err != nil {
		return nil, err
	}

	var timeout time.Duration
	if connType == MonType {
		timeout = MonSocketReadTimeout
	} else {
		timeout = CmdSocketReadTimeout
	}
	conn.SetReadDeadline(time.Now().Add(timeout))

	ch := NewFrameChannel(conn, conn, "BUS-"+connType)

	if connType == MonType {
		bc.monChannel = ch
		bc.monConn = conn
	} else {
		bc.cmdChannel = ch
		bc.cmdConn = conn
	}

	return ch, nil
}

func (bc *BUSConnector) doHandshake(ch *FrameChannel, connType string) error {
	log.Printf("(HS) starting HANDSHAKE on channel %s...", ch.GetName())

	// Set handshake timeout
	if connType == MonType && bc.monConn != nil {
		bc.monConn.SetReadDeadline(time.Now().Add(HandshakeTimeout))
	} else if bc.cmdConn != nil {
		bc.cmdConn.SetReadDeadline(time.Now().Add(HandshakeTimeout))
	}

	// STEP-1: wait for ACK from GW
	fr, err := ch.ReadFrames()
	if err != nil {
		return NewOWNAuthErrorWithCause("Handshake STEP-1 read error", err)
	}
	if fr != message2.FrameACK {
		return NewOWNAuthError(fmt.Sprintf("Could not open BUS-%s: no ACK at STEP-1, received: %s", connType, fr))
	}

	// STEP-2: send session request
	session := MonSession
	if connType == CmdType {
		session = CmdSession
	}
	if err := ch.SendFrame(session); err != nil {
		return NewOWNAuthErrorWithCause("Handshake STEP-2 send error", err)
	}

	fr, err = ch.ReadFrames()
	if err != nil {
		return NewOWNAuthErrorWithCause("Handshake STEP-2 read error", err)
	}

	if fr == message2.FrameNACK && connType == CmdType {
		// Try alt CMD session
		if err := ch.SendFrame(CmdSessionAlt); err != nil {
			return NewOWNAuthErrorWithCause("Handshake STEP-2 alt send error", err)
		}
		fr, err = ch.ReadFrames()
		if err != nil {
			return NewOWNAuthErrorWithCause("Handshake STEP-2 alt read error", err)
		}
	}

	if fr == message2.FrameACK {
		// NO_AUTH: no password required
		ch.HandshakeCompleted = true
		log.Println("(HS) NO_AUTH: GW has no pwd ==HANDSHAKE COMPLETED==")
	} else if matched, _ := regexp.MatchString(`^\*#\d+##$`, fr); matched {
		// OPEN_AUTH: nonce received
		if err := bc.doOPENHandshake(fr, ch); err != nil {
			return err
		}
		ch.HandshakeCompleted = true
	} else if fr == HMACSHA1 || fr == HMACSHA2 {
		// HMAC_AUTH
		if err := bc.doHMACHandshake(fr, ch); err != nil {
			return err
		}
		ch.HandshakeCompleted = true
	} else {
		return NewOWNAuthError(fmt.Sprintf("Cannot authenticate: unexpected answer: %s", fr))
	}

	// Reset read deadline for normal operation
	if connType == MonType && bc.monConn != nil {
		bc.monConn.SetReadDeadline(time.Now().Add(MonSocketReadTimeout))
	} else if bc.cmdConn != nil {
		bc.cmdConn.SetReadDeadline(time.Now().Add(CmdSocketReadTimeout))
	}

	return nil
}

func (bc *BUSConnector) doOPENHandshake(nonceFrame string, ch *FrameChannel) error {
	nonce := nonceFrame[2 : len(nonceFrame)-2]
	log.Printf("(HS) OPEN_AUTH: received nonce=%s", nonce)

	pwdEncoded, err := CalcOpenPass(bc.pwd, nonce)
	if err != nil {
		return NewOWNAuthError("Invalid gateway password. Password must contain only digits (OPEN_AUTH)")
	}

	pwdMessage := message2.FrameStartDim + pwdEncoded + message2.FrameEnd
	if err := ch.SendFrame(pwdMessage); err != nil {
		return NewOWNAuthErrorWithCause("OPEN_AUTH send error", err)
	}

	fr, err := ch.ReadFrames()
	if err != nil {
		return NewOWNAuthErrorWithCause("OPEN_AUTH read error", err)
	}
	if fr == message2.FrameACK {
		log.Println("(HS) OPEN_AUTH: pwd accepted ==HANDSHAKE COMPLETED==")
		return nil
	}
	return NewOWNAuthError("Password not accepted by gateway (OPEN_AUTH)")
}

func (bc *BUSConnector) doHMACHandshake(hmacType string, ch *FrameChannel) error {
	// STEP-3: send ACK, wait for Ra
	if err := ch.SendFrame(message2.FrameACK); err != nil {
		return NewOWNAuthErrorWithCause("HMAC STEP-3 send error", err)
	}

	fr, err := ch.ReadFrames()
	if err != nil {
		return NewOWNAuthErrorWithCause("HMAC STEP-3 read error", err)
	}

	pattern := regexp.MustCompile(`\*#(\d{80,128})##`)
	matches := pattern.FindStringSubmatch(fr)
	if len(matches) < 2 {
		return NewOWNAuthError(fmt.Sprintf("Handshake failed, no Ra received at HMAC STEP-3: %s", fr))
	}

	raDigits := matches[1]
	ra := DigitToHex(raDigits)
	rb := CalcHmacRb()
	a := "736F70653E"
	b := "636F70653E"
	kab := CalcSHA256(bc.pwd)
	hmacRaRbABKab := CalcSHA256(ra + rb + a + b + kab)

	// STEP-4: send Rb + HMAC
	hmacMessage := message2.FrameStartDim + HexToDigit(rb) + "*" + HexToDigit(hmacRaRbABKab) + message2.FrameEnd
	if err := ch.SendFrame(hmacMessage); err != nil {
		return NewOWNAuthErrorWithCause("HMAC STEP-4 send error", err)
	}

	fr, err = ch.ReadFrames()
	if err != nil {
		return NewOWNAuthErrorWithCause("HMAC STEP-4 read error", err)
	}

	if fr == message2.FrameNACK {
		return NewOWNAuthError("Password not accepted by gateway (HMAC)")
	}

	matches = pattern.FindStringSubmatch(fr)
	if len(matches) < 2 {
		return NewOWNAuthError(fmt.Sprintf("Handshake failed, invalid HMAC(Ra,Rb,Kab) at STEP-4: %s", fr))
	}

	hmacRaRbKab := DigitToHex(matches[1])
	if CalcSHA256(ra+rb+kab) == hmacRaRbKab {
		if err := ch.SendFrame(message2.FrameACK); err != nil {
			return NewOWNAuthErrorWithCause("HMAC final ACK send error", err)
		}
		log.Println("(HS) HMAC_AUTH: final ACK sent ==HANDSHAKE COMPLETED==")
		return nil
	}

	return NewOWNAuthError("Handshake failed, final HMAC(Ra,Rb,Kab) does not match")
}

// SendCommandSynch sends a command frame and waits for a response.
func (bc *BUSConnector) SendCommandSynch(frame string) (*Response, error) {
	if !bc.isCmdConnected {
		return nil, NewOWNError("CMD is not connected")
	}
	return bc.sendCmdAndReadResp(frame)
}

func (bc *BUSConnector) sendCmdAndReadResp(frame string) (*Response, error) {
	parsedReq, err := message2.Parse(frame)
	if err != nil {
		return nil, NewOWNErrorWithCause("Failed to parse request frame", err)
	}

	res := NewResponse(parsedReq)

	// Reset deadline for CMD
	if bc.cmdConn != nil {
		bc.cmdConn.SetReadDeadline(time.Now().Add(CmdSocketReadTimeout))
	}

	if err := bc.cmdChannel.SendFrame(frame); err != nil {
		// Try to reconnect CMD
		bc.isCmdConnected = false
		if reconnErr := bc.OpenCmdConn(); reconnErr != nil {
			return nil, NewOWNErrorWithCause("Cannot create new CMD connection", reconnErr)
		}
		// Retry
		if bc.cmdConn != nil {
			bc.cmdConn.SetReadDeadline(time.Now().Add(CmdSocketReadTimeout))
		}
		if err := bc.cmdChannel.SendFrame(frame); err != nil {
			return nil, NewOWNErrorWithCause("Failed to send on new CMD connection", err)
		}
	}

	bc.lastCmdFrameSentTs = time.Now().UnixMilli()

	for !res.HasFinalResponse() {
		fr, err := bc.cmdChannel.ReadFrames()
		if err != nil {
			return nil, NewOWNErrorWithCause("Error reading CMD response", err)
		}
		if fr == "" {
			return nil, NewOWNError("Received null frame while reading responses")
		}
		parsedResp, err := message2.Parse(fr)
		if err != nil {
			// Skip unsupported frames
			continue
		}
		res.AddResponse(parsedResp)
	}

	return res, nil
}

func (bc *BUSConnector) monReceiveLoop() {
	for {
		select {
		case <-bc.monStopChan:
			return
		default:
		}

		if bc.monConn != nil {
			bc.monConn.SetReadDeadline(time.Now().Add(MonSocketReadTimeout))
		}

		fr, err := bc.monChannel.ReadFrames()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Timeout - check if gateway is still reachable
				if bc.isCmdConnected {
					modelReq := message2.GatewayMgmtRequestModel()
					_, cmdErr := bc.sendCmdAndReadResp(modelReq.GetFrameValue())
					if cmdErr != nil {
						bc.handleMonDisconnect(NewOWNErrorWithCause("Gateway not reachable", cmdErr))
						return
					}
					continue
				}
			}

			select {
			case <-bc.monStopChan:
				return
			default:
			}

			bc.handleMonDisconnect(NewOWNErrorWithCause("MON read error", err))
			return
		}

		if fr == "" {
			bc.handleMonDisconnect(NewOWNError("MON readFrame() returned empty"))
			return
		}

		parsedMsg, err := message2.Parse(fr)
		if err != nil {
			log.Printf("##BUS-conn## Skipping frame: %s (%v)", fr, err)
			continue
		}
		if bc.listener != nil {
			bc.listener.OnMessage(parsedMsg)
		}
	}
}

func (bc *BUSConnector) startMonKeepalive() {
	bc.keepaliveDone = make(chan struct{})
	go func() {
		ticker := time.NewTicker(MonKeepaliveTimer)
		defer ticker.Stop()
		for {
			select {
			case <-bc.keepaliveDone:
				return
			case <-ticker.C:
				if bc.monChannel != nil {
					if err := bc.monChannel.SendFrame(message2.FrameACK); err != nil {
						log.Printf("##BUS-conn## Could not send MON keepalive: %v", err)
					}
				}
			}
		}
	}()
}

func (bc *BUSConnector) stopMonKeepalive() {
	if bc.keepaliveDone != nil {
		select {
		case <-bc.keepaliveDone:
		default:
			close(bc.keepaliveDone)
		}
	}
}

func (bc *BUSConnector) handleMonDisconnect(err error) {
	log.Printf("##BUS-conn## handleMonDisconnect: %v", err)
	bc.stopMonKeepalive()
	bc.isMonConnected = false
	if bc.monChannel != nil {
		bc.monChannel.Disconnect()
	}
	if bc.listener != nil {
		bc.listener.OnMonDisconnected(err)
	}
}

// Disconnect closes both MON and CMD connections.
func (bc *BUSConnector) Disconnect() {
	if bc.monStopChan != nil {
		select {
		case <-bc.monStopChan:
		default:
			close(bc.monStopChan)
		}
	}
	bc.stopMonKeepalive()
	bc.isCmdConnected = false
	bc.isMonConnected = false
	if bc.cmdChannel != nil {
		bc.cmdChannel.Disconnect()
	}
	if bc.monChannel != nil {
		bc.monChannel.Disconnect()
	}
	if bc.cmdConn != nil {
		bc.cmdConn.Close()
		bc.cmdConn = nil
	}
	if bc.monConn != nil {
		bc.monConn.Close()
		bc.monConn = nil
	}
	log.Println("##BUS-conn## CMD+MON connections CLOSED")
}
