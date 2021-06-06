package rest_errors

import (
	"encoding/json"
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
	FMessage string        `json:"message"`
	FStatus  int           `json:"status"`
	FError   string        `json:"error"`
	FCauses  []interface{} `json:"causes"`
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s; status: %d; error %s; causes[%v]",
		e.FMessage, e.FStatus, e.FError, e.FCauses)
}

func (e restErr) Message() string {
	return e.FMessage
}
func (e restErr) Status() int {
	return e.FStatus
}

func (e restErr) Causes() []interface{} {
	return e.FCauses
}

func NewRestError(msg string, status int, err string, causes []interface{}) error {
	return restErr{
		FMessage: msg,
		FStatus:  status,
		FError:   err,
		FCauses:  causes,
	}
}

func NewBadRequestError(msg string) RestErr {
	return restErr{
		FMessage: msg,
		FStatus:  http.StatusBadRequest,
		FError:   "bad request",
	}
}

func NewAuthorizationError(msg string) RestErr {
	return restErr{
		FMessage: msg,
		FStatus:  http.StatusUnauthorized,
		FError:   "unauthorized",
	}
}

func NewNotFoundError(msg string) RestErr {
	return restErr{
		FMessage: msg,
		FStatus:  http.StatusNotFound,
		FError:   "not found",
	}
}

func NewInternalServerError(msg string, err error) RestErr {
	return restErr{
		FMessage: msg,
		FStatus:  http.StatusInternalServerError,
		FError:   "internal server error",
		FCauses:  []interface{}{err},
	}
}

func ResponseJson(w http.ResponseWriter, httpStatus int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(body)
}

func ResponseError(w http.ResponseWriter, err RestErr) {
	ResponseJson(w, err.Status(), err)
}
