package communication

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/jensvandenreyt/openwebnet4go/message"
	"github.com/rs/zerolog/log"
)

// FrameChannel wraps input/output streams to send and receive frames from an OpenWebNet gateway.
type FrameChannel struct {
	reader             *bufio.Reader
	writer             io.Writer
	name               string
	HandshakeCompleted bool
	readQueue          []string
}

// NewFrameChannel creates a new FrameChannel.
func NewFrameChannel(r io.Reader, w io.Writer, name string) *FrameChannel {
	return &FrameChannel{
		reader:    bufio.NewReader(r),
		writer:    w,
		name:      name,
		readQueue: make([]string, 0),
	}
}

// GetName returns the name of this channel.
func (fc *FrameChannel) GetName() string {
	return fc.name
}

// SendFrame sends a frame on the channel.
func (fc *FrameChannel) SendFrame(frame string) error {
	if fc.writer == nil {
		return fmt.Errorf("cannot sendFrame, writer is nil")
	}
	_, err := fc.writer.Write([]byte(frame))
	if err != nil {
		return err
	}
	log.Trace().Msgf("-FC-%s -------> %s", fc.name, frame)
	return nil
}

// ReadFrames returns the first frame from the receiving queue.
// If the queue is empty, reads from the underlying reader until a frame delimiter (##) is found.
// Returns empty string if end of stream is reached.
func (fc *FrameChannel) ReadFrames() (string, error) {
	if len(fc.readQueue) > 0 {
		frame := fc.readQueue[0]
		fc.readQueue = fc.readQueue[1:]
		return frame, nil
	}

	// Read until we get a frame ending with ##
	buf, err := fc.readUntilDelimiter()
	if err != nil {
		return "", err
	}
	if len(buf) == 0 {
		return "", nil
	}

	longFrame := string(buf)

	if strings.Contains(longFrame, message.FrameEnd) {
		frames := strings.Split(longFrame, message.FrameEnd)
		for _, f := range frames {
			if f != "" {
				fc.readQueue = append(fc.readQueue, f+message.FrameEnd)
			}
		}
		if len(fc.readQueue) > 0 {
			frame := fc.readQueue[0]
			fc.readQueue = fc.readQueue[1:]
			log.Trace().Msgf("-FC-%s <------- %s", fc.name, frame)
			return frame, nil
		}
	}

	return "", fmt.Errorf("no delimiter found on stream: %s", longFrame)
}

// readUntilDelimiter reads from the reader until the delimiter (##) is found.
func (fc *FrameChannel) readUntilDelimiter() ([]byte, error) {
	var buf []byte
	hashFound := false

	for {
		b, err := fc.reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return buf, nil
			}
			return nil, err
		}

		buf = append(buf, b)
		c := rune(b)

		if c == '#' && !hashFound {
			hashFound = true
		} else if c == '#' && hashFound {
			// Found ##, frame terminated
			break
		} else {
			hashFound = false
		}
	}

	return buf, nil
}

// Disconnect closes the channel.
func (fc *FrameChannel) Disconnect() {
	if closer, ok := fc.writer.(io.Closer); ok {
		closer.Close()
	}
	fc.writer = nil
	// bufio.Reader does not implement io.Closer, so nothing to close here.
	// The underlying connection is closed separately.
	log.Trace().Msgf("-FC-%s in/out streams CLOSED", fc.name)
}
