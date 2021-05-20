package rest_errors

import (
	"fmt"
	"net/http"
)

// client must use interface
type RestErr interface {
	Message() string
	Status() int
	Error() string
	Causes() []interface{}
}
type restErr struct {
	message string        `json:"message"`
	status  int           `json:"status"`
	error   string        `json:"error"`
	causes  []interface{} `json:"causes"`
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s; status: %d; error %s; causes[%v]", e.message, e.status, e.error, e.causes)
}

func (e restErr) Message() string {
	return e.message
}
func (e restErr) Status() int {
	return e.status
}

func (e restErr) Causes() []interface{} {
	return e.causes
}

func NewRestError(msg string, status int, err string, causes []interface{}) error {
	return restErr{
		message: msg,
		status:  status,
		error:   err,
		causes:  causes,
	}
}

func NewBadRequestError(msg string) RestErr {
	return restErr{
		message: msg,
		status:  http.StatusBadRequest,
		error:   "bad request",
	}
}

func NewAuthorizationError(msg string) RestErr {
	return restErr{
		message: msg,
		status:  http.StatusUnauthorized,
		error:   "unauthorized",
	}
}

func NewNotFoundError(msg string) RestErr {
	return restErr{
		message: msg,
		status:  http.StatusNotFound,
		error:   "not found",
	}
}

func NewInternalServerError(msg string, err error) RestErr {
	return restErr{
		message: msg,
		status:  http.StatusInternalServerError,
		error:   "internal server error",
		causes:  []interface{}{err},
	}
}
