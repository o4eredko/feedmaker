// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// StartStopper is an autogenerated mock type for the StartStopper type
type StartStopper struct {
	mock.Mock
}

// Start provides a mock function with given fields:
func (_m *StartStopper) Start() {
	_m.Called()
}

// Stop provides a mock function with given fields:
func (_m *StartStopper) Stop() context.Context {
	ret := _m.Called()

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}
