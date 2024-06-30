package api

import "net/http"

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func ErrorResourceNotFound(res string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err: res + " resource not found",
	}
}

func ErrorBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err: "invalid JSON request",
	}
}

func ErrorInvalidID() Error{
	return Error{
		Code: http.StatusBadRequest,
		Err: "invalid id",
	}
}