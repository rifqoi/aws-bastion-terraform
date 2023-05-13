package common

import "fmt"

// Safer way to create enum in golang
// Reference from https://threedots.tech/post/safer-enums-in-go/
type ErrorType struct {
	s string
}

var (
	ErrorTypeUnknown       = ErrorType{"unknown"}
	ErrorTypeNotFound      = ErrorType{"not-found"}
	ErrorTypeServerError   = ErrorType{"server-error"}
	ErrorTypeAuthorization = ErrorType{"unauthorized"}
	ErrorTypeInvalidInput  = ErrorType{"invalid-input"}
)

// Create a custom error struct that implement error interface
type Error struct {
	origin  error
	msg     string
	errType ErrorType
}

// Wrapping error to manage fine grained error handling
func WrapErrorf(orig error, errType ErrorType, format string, a ...any) error {
	return &Error{
		origin:  orig,
		errType: errType,
		msg:     fmt.Sprintf(format, a...),
	}
}

func NewErrorf(errType ErrorType, format string, a ...any) error {
	return WrapErrorf(nil, errType, format, a...)

}

func (e *Error) Error() string {
	if e.origin != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.origin)
	}

	return e.msg
}

func (e *Error) Unwrap() error {
	return e.origin
}

func (e *Error) Type() ErrorType {
	return e.errType
}
