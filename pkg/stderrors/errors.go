package stderrors

import (
	"net/http"
	"strings"
)

type WebError interface {
	error
	GetPublicMessage() string
	GetStatusCode() int
}

type Error struct {
	Code    int
	Message string

	// Used to store information about error that would appear only in application logs, but not in server responses
	InternalMessage string
}

func (e *Error) Error() string {
	return strings.Join([]string{e.Message, e.InternalMessage}, " - ")
}

func (e *Error) GetPublicMessage() string {
	return e.Message
}

func (e *Error) GetStatusCode() int {
	return e.Code
}

func BadRequest(message string) *Error {
	return &Error{Code: http.StatusBadRequest, Message: message}
}

func Unauthorized(message string) *Error {
	return &Error{Code: http.StatusUnauthorized, Message: message}
}

func NotFound(message string) *Error {
	return &Error{Code: http.StatusNotFound, Message: message}
}

func UnprocessableEntity(message string) *Error {
	return &Error{Code: http.StatusUnprocessableEntity, Message: message}
}
