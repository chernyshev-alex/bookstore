package rest_errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInternalServerError(t *testing.T) {
	err := NewInternalServerError("msg", errors.New("cause"))
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, 1, len(err.Causes()))
}

func TestNewBadRequestError(t *testing.T) {
	err := NewBadRequestError("msg")
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("msg")
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}
