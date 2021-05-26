package mocks

import (
	"net/http"

	rest_errors "github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
	mock "github.com/stretchr/testify/mock"
)

type MockOAuth struct {
	mock.Mock
}

func (m *MockOAuth) AuthenticateRequest(rq *http.Request) rest_errors.RestErr {
	args := m.Called(rq)
	return args.Get(0).(rest_errors.RestErr)
}

func (m *MockOAuth) IsPublic(rq *http.Request) bool {
	args := m.Called(rq)
	return args.Bool(0)
}

func (m *MockOAuth) GetCallerId(rq *http.Request) int64 {
	args := m.Called(rq)
	return int64(args.Int(0))
}
func (m *MockOAuth) GetClientId(rq *http.Request) int64 {
	args := m.Called(rq)
	return int64(args.Int(0))
}
