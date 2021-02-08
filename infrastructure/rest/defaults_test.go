package rest_test

import (
	"errors"

	"go-feedmaker/interactor/mocks"
)

type (
	handlerFields struct {
		feeds *mocks.FeedInteractor
	}
)

var (
	defaultSentinel = "foo, bar, baz"
	defaultTestErr  = errors.New("default test error")
)

func defaultHandlerFields() *handlerFields {
	return &handlerFields{
		feeds: new(mocks.FeedInteractor),
	}
}
