package communication

import (
	"fmt"
	"sync"

	"github.com/jensvandenreyt/openwebnet4go/message"
)

// Response represents a response to an OpenWebNet request sent to the gateway.
type Response struct {
	requestMessage message.OpenMessage
	finalResponse  message.OpenMessage
	responses      []message.OpenMessage
	IsSuccess      bool
	mu             sync.Mutex
	done           chan struct{}
}

// NewResponse creates a new Response associated with the request message.
func NewResponse(request message.OpenMessage) *Response {
	return &Response{
		requestMessage: request,
		responses:      make([]message.OpenMessage, 0),
		done:           make(chan struct{}),
	}
}

// GetRequest returns the initial request message.
func (r *Response) GetRequest() message.OpenMessage {
	return r.requestMessage
}

// GetResponseMessages returns the list of OpenMessages received as response.
func (r *Response) GetResponseMessages() []message.OpenMessage {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]message.OpenMessage, len(r.responses))
	copy(result, r.responses)
	return result
}

// GetFirstBaseOpenMessage returns the first of BaseOpenMessage received as response.
func (r *Response) GetFirstBaseOpenMessage() *message.BaseOpenMessage {
	messages := r.GetResponseMessages()
	for _, openMessage := range messages {
		if bom, ok := openMessage.(*message.BaseOpenMessage); ok {
			return bom
		}
	}
	return nil
}

// GetFinalResponse returns the last OpenMessage that finalised this response.
func (r *Response) GetFinalResponse() message.OpenMessage {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.finalResponse
}

// String returns a string representation of the response.
func (r *Response) String() string {
	return fmt.Sprintf("{REQ=%s|RESP=%v}", r.requestMessage.String(), r.responses)
}

// AddResponse adds a new message received as response.
func (r *Response) AddResponse(msg message.OpenMessage) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.responses = append(r.responses, msg)
	if msg.IsACK() || msg.IsNACK() {
		r.finalResponse = msg
		if msg.IsACK() {
			r.IsSuccess = true
		}
		// Signal that the response is complete
		select {
		case <-r.done:
			// already closed
		default:
			close(r.done)
		}
	}
}

// HasFinalResponse returns true if an ACK/NACK has been received.
func (r *Response) HasFinalResponse() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.finalResponse != nil
}

// WaitResponse blocks until a final response is received.
func (r *Response) WaitResponse() {
	r.mu.Lock()
	if r.finalResponse != nil {
		r.mu.Unlock()
		return
	}
	r.mu.Unlock()
	<-r.done
}

// ResponseReady signals that the response is ready (called externally if needed).
func (r *Response) ResponseReady() {
	select {
	case <-r.done:
		// already closed
	default:
		close(r.done)
	}
}
