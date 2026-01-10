package api

import "net/http"

type AppError struct {
	Err    error
	Msg    string
	Status int
}

func (e *AppError) Error() string {
	return e.Msg
}

func NewError(err error, status int, msg string) *AppError {
	return &AppError{
		Err:    err,
		Status: status,
		Msg:    msg,
	}
}

func ErrInternal(err error) *AppError {
	return NewError(err, http.StatusInternalServerError, "internal_error")
}

func ErrBadRequest(msg string) *AppError {
	return NewError(nil, http.StatusBadRequest, msg)
}

func ErrUnauthorized(msg string) *AppError {
	return NewError(nil, http.StatusUnauthorized, msg)
}

func ErrNotFound(msg string) *AppError {
	return NewError(nil, http.StatusNotFound, msg)
}
