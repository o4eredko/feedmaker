// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	broadcaster "go-feedmaker/infrastructure/rest/broadcaster"

	mock "github.com/stretchr/testify/mock"
)

// Recipient is an autogenerated mock type for the Recipient type
type Recipient struct {
	mock.Mock
}

// OnCloseHook provides a mock function with given fields: hook
func (_m *Recipient) OnCloseHook(hook broadcaster.CloseHook) {
	_m.Called(hook)
}

// Send provides a mock function with given fields: _a0
func (_m *Recipient) Send(_a0 []byte) {
	_m.Called(_a0)
}

// Start provides a mock function with given fields:
func (_m *Recipient) Start() {
	_m.Called()
}

// Stop provides a mock function with given fields:
func (_m *Recipient) Stop() {
	_m.Called()
}
