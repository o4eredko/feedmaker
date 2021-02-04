// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"
	scheduler "go-feedmaker/infrastructure/scheduler"

	mock "github.com/stretchr/testify/mock"
)

// TaskScheduler is an autogenerated mock type for the TaskScheduler type
type TaskScheduler struct {
	mock.Mock
}

// AddTask provides a mock function with given fields: taskID, task
func (_m *TaskScheduler) AddTask(taskID scheduler.TaskID, task scheduler.Task) error {
	ret := _m.Called(taskID, task)

	var r0 error
	if rf, ok := ret.Get(0).(func(scheduler.TaskID, scheduler.Task) error); ok {
		r0 = rf(taskID, task)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveTask provides a mock function with given fields: taskID
func (_m *TaskScheduler) RemoveTask(taskID scheduler.TaskID) error {
	ret := _m.Called(taskID)

	var r0 error
	if rf, ok := ret.Get(0).(func(scheduler.TaskID) error); ok {
		r0 = rf(taskID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields:
func (_m *TaskScheduler) Start() {
	_m.Called()
}

// Stop provides a mock function with given fields:
func (_m *TaskScheduler) Stop() context.Context {
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
