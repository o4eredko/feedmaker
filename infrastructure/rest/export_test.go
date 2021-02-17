package rest

import (
	"net/http"

	"go-feedmaker/infrastructure/scheduler"
	"go-feedmaker/interactor"
)

type (
	ScheduleTaskIn = scheduleTaskIn
	ScheduleOut    = scheduleOut
)

func MakeSchedulesOut(schedules map[scheduler.TaskID]*scheduler.Schedule) map[scheduler.TaskID]*scheduleOut {
	return makeSchedulesOut(schedules)
}

func (h *handler) Feeds() interactor.FeedInteractor {
	return h.feeds
}

func (s *APIServer) Server() *http.Server {
	return s.server
}

func (h *wsHandler) Upgrader() Upgrader {
	return h.upgrader
}

func (h *wsHandler) Broadcaster() Broadcaster {
	return h.broadcaster
}
