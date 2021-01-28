package repository

import "go-feedmaker/interactor"

type (
	defaultFactory struct{}
)

func NewDefaultFactory() interactor.FeedFactory {
	return nil
}

func NewYandexFactory() interactor.FeedFactory {
	return nil
}
