package message

import "fmt"

// FrameError is used when a problem with an OpenWebNet frame occurred.
type FrameError struct {
	Message string
	Cause   error
}

func (e *FrameError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *FrameError) Unwrap() error {
	return e.Cause
}

// MalformedFrameError is used when an OpenWebNet frame is malformed.
type MalformedFrameError struct {
	FrameError
}

// UnsupportedFrameError is used when an OpenWebNet frame is not supported by the library.
type UnsupportedFrameError struct {
	FrameError
}

func NewFrameError(msg string) *FrameError {
	return &FrameError{Message: msg}
}

func NewFrameErrorWithCause(msg string, cause error) *FrameError {
	return &FrameError{Message: msg, Cause: cause}
}

func NewMalformedFrameError(msg string) *MalformedFrameError {
	return &MalformedFrameError{FrameError{Message: msg}}
}

func NewUnsupportedFrameError(msg string) *UnsupportedFrameError {
	return &UnsupportedFrameError{FrameError{Message: msg}}
}
