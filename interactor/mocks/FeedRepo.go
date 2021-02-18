// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"
	entity "go-feedmaker/entity"

	interactor "go-feedmaker/interactor"

	mock "github.com/stretchr/testify/mock"
)

// FeedRepo is an autogenerated mock type for the FeedRepo type
type FeedRepo struct {
	mock.Mock
}

// CancelGeneration provides a mock function with given fields: ctx, id
func (_m *FeedRepo) CancelGeneration(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetFactoryByGenerationType provides a mock function with given fields: generationType
func (_m *FeedRepo) GetFactoryByGenerationType(generationType string) (interactor.FeedFactory, error) {
	ret := _m.Called(generationType)

	var r0 interactor.FeedFactory
	if rf, ok := ret.Get(0).(func(string) interactor.FeedFactory); ok {
		r0 = rf(generationType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interactor.FeedFactory)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(generationType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGeneration provides a mock function with given fields: ctx, generationID
func (_m *FeedRepo) GetGeneration(ctx context.Context, generationID string) (*entity.Generation, error) {
	ret := _m.Called(ctx, generationID)

	var r0 *entity.Generation
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.Generation); ok {
		r0 = rf(ctx, generationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Generation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, generationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsAllowedType provides a mock function with given fields: generationType
func (_m *FeedRepo) IsAllowedType(generationType string) bool {
	ret := _m.Called(generationType)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(generationType)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ListAllowedTypes provides a mock function with given fields:
func (_m *FeedRepo) ListAllowedTypes() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// ListGenerations provides a mock function with given fields: ctx
func (_m *FeedRepo) ListGenerations(ctx context.Context) ([]*entity.Generation, error) {
	ret := _m.Called(ctx)

	var r0 []*entity.Generation
	if rf, ok := ret.Get(0).(func(context.Context) []*entity.Generation); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.Generation)
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

// OnGenerationCanceled provides a mock function with given fields: ctx, id, callback
func (_m *FeedRepo) OnGenerationCanceled(ctx context.Context, id string, callback func()) error {
	ret := _m.Called(ctx, id, callback)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, func()) error); ok {
		r0 = rf(ctx, id, callback)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OnGenerationsUpdated provides a mock function with given fields: ctx, callback
func (_m *FeedRepo) OnGenerationsUpdated(ctx context.Context, callback func(*entity.Generation)) error {
	ret := _m.Called(ctx, callback)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(*entity.Generation)) error); ok {
		r0 = rf(ctx, callback)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StoreGeneration provides a mock function with given fields: ctx, generation
func (_m *FeedRepo) StoreGeneration(ctx context.Context, generation *entity.Generation) error {
	ret := _m.Called(ctx, generation)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Generation) error); ok {
		r0 = rf(ctx, generation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateGenerationState provides a mock function with given fields: ctx, generation
func (_m *FeedRepo) UpdateGenerationState(ctx context.Context, generation *entity.Generation) error {
	ret := _m.Called(ctx, generation)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Generation) error); ok {
		r0 = rf(ctx, generation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
