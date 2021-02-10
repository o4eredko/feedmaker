package scheduler_test

import (
	"errors"

	"github.com/robfig/cron/v3"

	"go-feedmaker/infrastructure/scheduler"
	"go-feedmaker/infrastructure/scheduler/mocks"
)

var (
	defaultErr     = errors.New("default error")
	defaultEntryID = cron.EntryID(42)
	defaultTaskID  = scheduler.TaskID("foobar")
	defaultTask    = getMockedTask()
)

func getMockedTask() *mocks.Task {
	mockedTask := new(mocks.Task)
	mockedTask.
		On("Schedule").Return(new(mocks.Schedule)).
		On("Cmd").Return(new(mocks.Runner))
	return mockedTask
}
