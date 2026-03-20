package message

// AckOpenMessage represents ACK, NACK, and BUSY-NACK messages.
type AckOpenMessage struct {
	BaseMessage
}

// Pre-constructed singleton instances.
var (
	ACK      = &AckOpenMessage{BaseMessage{FrameValue: FrameACK}}
	NACK     = &AckOpenMessage{BaseMessage{FrameValue: FrameNACK}}
	BusyNACK = &AckOpenMessage{BaseMessage{FrameValue: FrameBusyNACK}}
)

func NewAckOpenMessage(frame string) *AckOpenMessage {
	return &AckOpenMessage{BaseMessage{FrameValue: frame}}
}

func (m *AckOpenMessage) IsCommand() bool {
	return false
}

func (m *AckOpenMessage) ToStringVerbose() string {
	return "<" + m.FrameValue + ">"
}
