// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// WSHandler is an autogenerated mock type for the WSHandler type
type WSHandler struct {
	mock.Mock
}

// ServeWS provides a mock function with given fields: w, r
func (_m *WSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
