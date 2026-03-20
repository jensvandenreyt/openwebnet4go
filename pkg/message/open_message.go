package message

// Frame constants for OpenWebNet protocol.
const (
	FrameACK          = "*#*1##"
	FrameNACK         = "*#*0##"
	FrameBusyNACK     = "*#*6##"
	FrameACKNACKStart = "*#*"
	FrameStart        = "*"
	FrameStartDim     = "*#"
	FrameEnd          = "##"
)

// OpenMessage is the interface for all OpenWebNet messages.
type OpenMessage interface {
	// GetFrameValue returns the raw frame value.
	GetFrameValue() string
	// IsCommand returns true if this is a command frame (*WHO..).
	IsCommand() bool
	// IsACK returns true if this is an ACK frame (*#*1##).
	IsACK() bool
	// IsNACK returns true if this is a NACK frame (*#*0##).
	IsNACK() bool
	// IsBusyNACK returns true if this is a BUSY_NACK frame (*#*6##).
	IsBusyNACK() bool
	// ToStringVerbose returns a verbose representation of the message.
	ToStringVerbose() string
	// String returns a short string representation of the message.
	String() string
}

// BaseMessage contains common fields shared by all message types.
type BaseMessage struct {
	FrameValue string
}

func (m *BaseMessage) GetFrameValue() string {
	return m.FrameValue
}

func (m *BaseMessage) IsACK() bool {
	return m.FrameValue == FrameACK
}

func (m *BaseMessage) IsNACK() bool {
	return m.FrameValue == FrameNACK
}

func (m *BaseMessage) IsBusyNACK() bool {
	return m.FrameValue == FrameBusyNACK
}

func (m *BaseMessage) String() string {
	return "<" + m.FrameValue + ">"
}
