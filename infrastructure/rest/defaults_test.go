package rest_test

import (
	"errors"

	restMocks "go-feedmaker/infrastructure/rest/mocks"
	interactorMocks "go-feedmaker/interactor/mocks"
)

type (
	handlerFields struct {
		feeds     *interactorMocks.FeedInteractor
		scheduler *restMocks.Scheduler
	}
)

var (
	defaultSentinel = "foo, bar, baz"
	defaultTestErr  = errors.New("default test error")
)

func defaultHandlerFields() *handlerFields {
	return &handlerFields{
		feeds:     new(interactorMocks.FeedInteractor),
		scheduler: new(restMocks.Scheduler),
	}
}
