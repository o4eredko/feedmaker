// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	interactor "go-feedmaker/interactor"

	mock "github.com/stretchr/testify/mock"
)

// FeedFactory is an autogenerated mock type for the FeedFactory type
type FeedFactory struct {
	mock.Mock
}

// CreateDataFetcher provides a mock function with given fields:
func (_m *FeedFactory) CreateDataFetcher() interactor.DataFetcher {
	ret := _m.Called()

	var r0 interactor.DataFetcher
	if rf, ok := ret.Get(0).(func() interactor.DataFetcher); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interactor.DataFetcher)
		}
	}

	return r0
}

// CreateFileFormatter provides a mock function with given fields: dataStream
func (_m *FeedFactory) CreateFileFormatter(dataStream <-chan []string) interactor.FileFormatter {
	ret := _m.Called(dataStream)

	var r0 interactor.FileFormatter
	if rf, ok := ret.Get(0).(func(<-chan []string) interactor.FileFormatter); ok {
		r0 = rf(dataStream)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interactor.FileFormatter)
		}
	}

	return r0
}
