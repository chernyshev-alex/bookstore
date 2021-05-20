package http_utils

import (
	"encoding/json"
	"net/http"

	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
)

func ResponseJson(w http.ResponseWriter, httpStatus int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(body)
}

func ResponseError(w http.ResponseWriter, err rest_errors.RestErr) {
	ResponseJson(w, err.Status(), err)
}
