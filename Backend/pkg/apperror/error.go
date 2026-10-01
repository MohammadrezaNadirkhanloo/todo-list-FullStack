package apperror

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    Code
	Message string
	Op      string
	Err     error
}

func (e *Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
}

func (e *Error) Unwrap() error { return e.Err }

func New(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func Wrap(err error, code Code, message string) error {
	if err == nil {
		return nil
	}
	return &Error{Code: code, Message: message, Err: err}
}

func WithOp(err error, op string) error {
	if err == nil {
		return nil
	}
	var appErr *Error
	if errors.As(err, &appErr) {
		if appErr.Op == "" {
			appErr.Op = op
		}
		return err
	}
	return &Error{Code: CodeInternal, Message: msgInternal, Op: op, Err: err}
}

func CodeOf(err error) Code {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return CodeInternal
}

const msgInternal = "an internal server error occurred."

func UserMessage(err error) string {
	var appErr *Error
	if errors.As(err, &appErr) && appErr.Message != "" {
		return appErr.Message
	}
	return msgInternal
}

func NotFound(entity string) *Error {
	return &Error{Code: CodeNotFound, Message: fmt.Sprintf("%s not found.", entity)}
}
func Conflict(message string) *Error {
	return &Error{Code: CodeConflict, Message: message}
}
func Unauthorized(message string) *Error {
	return &Error{Code: CodeUnauthorized, Message: message}
}
func Forbidden() *Error {
	return &Error{Code: CodeForbidden, Message: "you do not have permission to perform this operation."}
}
func InvalidInput(message string) *Error {
	return &Error{Code: CodeInvalidInput, Message: message}
}

func Internal(err error) error {
	return Wrap(err, CodeInternal, msgInternal)
}
