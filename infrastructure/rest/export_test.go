package rest

import (
	"net/http"

	"go-feedmaker/infrastructure/scheduler"
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

func MakeSchedulesOut(schedules map[scheduler.TaskID]*scheduler.Schedule) []*scheduleOut {
	return makeSchedulesOut(schedules)
}
