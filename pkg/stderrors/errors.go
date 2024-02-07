package stderrors

import "net/http"

type WebError interface {
	error
	GetStatusCode() int
}

type Error struct {
	Code    int
	Message string
}

func (n *Error) Error() string {
	return n.Message
}

func (n *Error) GetStatusCode() int {
	return http.StatusNotFound
}

func BadRequest(message string) *Error {
	return &Error{Code: http.StatusBadRequest, Message: message}
}

func NotFound(message string) *Error {
	return &Error{Code: http.StatusNotFound, Message: message}
}

func UnprocessableEntity(message string) *Error {
	return &Error{Code: http.StatusUnprocessableEntity, Message: message}
}
