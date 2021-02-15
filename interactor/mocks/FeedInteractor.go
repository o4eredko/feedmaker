// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"
	entity "go-feedmaker/entity"

	mock "github.com/stretchr/testify/mock"
)

// FeedInteractor is an autogenerated mock type for the FeedInteractor type
type FeedInteractor struct {
	mock.Mock
}

// CancelGeneration provides a mock function with given fields: ctx, id
func (_m *FeedInteractor) CancelGeneration(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GenerateFeed provides a mock function with given fields: ctx, generationType
func (_m *FeedInteractor) GenerateFeed(ctx context.Context, generationType string) error {
	ret := _m.Called(ctx, generationType)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, generationType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListGenerationTypes provides a mock function with given fields: ctx
func (_m *FeedInteractor) ListGenerationTypes(ctx context.Context) (interface{}, error) {
	ret := _m.Called(ctx)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context) interface{}); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListGenerations provides a mock function with given fields: ctx
func (_m *FeedInteractor) ListGenerations(ctx context.Context) (interface{}, error) {
	ret := _m.Called(ctx)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context) interface{}); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WatchGenerationsProgress provides a mock function with given fields: ctx, outStream
func (_m *FeedInteractor) WatchGenerationsProgress(ctx context.Context, outStream chan<- *entity.Generation) error {
	ret := _m.Called(ctx, outStream)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, chan<- *entity.Generation) error); ok {
		r0 = rf(ctx, outStream)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
