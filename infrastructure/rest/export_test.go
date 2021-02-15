package rest

import (
	"net/http"

	"go-feedmaker/interactor"
)

type (
	ScheduleTaskIn = scheduleTaskIn
)

func (h *handler) Feeds() interactor.FeedInteractor {
	return h.feeds
}

func (s *APIServer) Server() *http.Server {
	return s.server
}
