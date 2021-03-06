// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"

	rest_errors "github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/rest_errors"
)

// OAuthInterface is an autogenerated mock type for the OAuthInterface type
type OAuthInterface struct {
	mock.Mock
}

// AuthenticateRequest provides a mock function with given fields: _a0
func (_m *OAuthInterface) AuthenticateRequest(_a0 *http.Request) rest_errors.RestErr {
	ret := _m.Called(_a0)

	var r0 rest_errors.RestErr
	if rf, ok := ret.Get(0).(func(*http.Request) rest_errors.RestErr); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(rest_errors.RestErr)
		}
	}

	return r0
}

// GetCallerId provides a mock function with given fields: _a0
func (_m *OAuthInterface) GetCallerId(_a0 *http.Request) int64 {
	ret := _m.Called(_a0)

	var r0 int64
	if rf, ok := ret.Get(0).(func(*http.Request) int64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// GetClientId provides a mock function with given fields: _a0
func (_m *OAuthInterface) GetClientId(_a0 *http.Request) int64 {
	ret := _m.Called(_a0)

	var r0 int64
	if rf, ok := ret.Get(0).(func(*http.Request) int64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// IsPublic provides a mock function with given fields: _a0
func (_m *OAuthInterface) IsPublic(_a0 *http.Request) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*http.Request) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
