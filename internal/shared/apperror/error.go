package apperror

import (
	"errors"
	"fmt"
)

type Code string

const (
	CodeBadRequest Code = "BAD_REQUEST"
	CodeNotFound   Code = "NOT_FOUND"
	CodeInternal   Code = "INTERNAL"
)

type Error struct {
	Code    Code
	Message string
	Errors  any
	Err     error
}

func New(code Code, message string, errors any, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Errors:  errors,
		Err:     err,
	}
}

func BadRequest(message string, errors any) *Error {
	return New(CodeBadRequest, message, errors, nil)
}

func NotFound(message string) *Error {
	return New(CodeNotFound, message, nil, nil)
}

func Internal(err error) *Error {
	return New(CodeInternal, "Internal server error", nil, err)
}

func As(err error) (*Error, bool) {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr, true
	}

	return nil, false
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}

	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}
