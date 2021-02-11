package rest_test

import (
	"errors"
	"time"

	"go-feedmaker/infrastructure/rest"
	restMocks "go-feedmaker/infrastructure/rest/mocks"
	"go-feedmaker/infrastructure/scheduler"
	"go-feedmaker/infrastructure/scheduler/task"
	interactorMocks "go-feedmaker/interactor/mocks"
)

type (
	handlerFields struct {
		feeds     *interactorMocks.FeedInteractor
		scheduler *restMocks.Scheduler
	}
)

var (
	defaultSentinel   = "foo, bar, baz"
	defaultTestErr    = errors.New("default test error")
	defaultScheduleIn = rest.ScheduleTaskIn{
		StartTimestamp: time.Now().UTC().Add(time.Hour * 13),
		DelayInterval:  time.Minute * 42,
	}
	defaultTaskSchedules = map[scheduler.TaskID]*task.Schedule{
		"foobar": {}, "spam": {}, "ham": {}, "eggs": {}, "0xDEADBEEF": {},
	}
)

func defaultHandlerFields() *handlerFields {
	return &handlerFields{
		feeds:     new(interactorMocks.FeedInteractor),
		scheduler: new(restMocks.Scheduler),
	}
}
