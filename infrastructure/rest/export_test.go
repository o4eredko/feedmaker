package rest

import (
	"net/http"

	"go-feedmaker/interactor"
)

func (h *handler) Feeds() interactor.FeedInteractor {
	return h.feeds
}

func (a *API) Server() *http.Server {
	return a.server
}
