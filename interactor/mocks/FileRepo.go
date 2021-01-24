// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	io "io"

	mock "github.com/stretchr/testify/mock"
)

// FileRepo is an autogenerated mock type for the FileRepo type
type FileRepo struct {
	mock.Mock
}

// Upload provides a mock function with given fields: ctx, file
func (_m *FileRepo) Upload(ctx context.Context, file io.Reader) error {
	ret := _m.Called(ctx, file)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, io.Reader) error); ok {
		r0 = rf(ctx, file)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
