package communication

import "fmt"

// OWNError represents an OpenWebNet error.
type OWNError struct {
	Message string
	Cause   error
}

func (e *OWNError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *OWNError) Unwrap() error {
	return e.Cause
}

func NewOWNError(msg string) *OWNError {
	return &OWNError{Message: msg}
}

func NewOWNErrorWithCause(msg string, cause error) *OWNError {
	return &OWNError{Message: msg, Cause: cause}
}

// OWNAuthError represents an authentication error.
type OWNAuthError struct {
	OWNError
}

func NewOWNAuthError(msg string) *OWNAuthError {
	return &OWNAuthError{OWNError{Message: msg}}
}

func NewOWNAuthErrorWithCause(msg string, cause error) *OWNAuthError {
	return &OWNAuthError{OWNError{Message: msg, Cause: cause}}
}
