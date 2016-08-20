package errors

import (
	"fmt"
)

func (e *Errors) StatusCode() int32 {
	return e.Code
}

func (e *Errors) Error() string {
	return e.Message
}

func New(code int32, msg string) error {
	return &Errors{
		Code:    code,
		Message: msg,
	}
}

func Errorhf(code int32, format string, args ...interface{}) error {
	return New(code, fmt.Sprintf(format, args...))
}

func BadRequest(message string) error {
	return New(400, message)
}

func Unauthorized(message string) error {
	return New(401, message)
}

func Forbidden(message string) error {
	return New(403, message)
}

func NotFound(message string) error {
	return New(404, message)
}

func InternalServerError(message string) error {
	return New(405, message)
}
