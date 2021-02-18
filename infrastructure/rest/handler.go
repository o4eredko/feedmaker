package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-feedmaker/infrastructure/scheduler"
	"go-feedmaker/interactor"
)

type (
	handler struct {
		feeds     interactor.FeedInteractor
		scheduler Scheduler
	}

	Scheduler interface {
		ScheduleTask(taskID scheduler.TaskID, task *scheduler.Task) error
		RemoveTask(taskID scheduler.TaskID) error
		ListSchedules() (map[scheduler.TaskID]*scheduler.Schedule, error)
	}
)

func NewHandler(feeds interactor.FeedInteractor, scheduler Scheduler) *handler {
	return &handler{
		feeds:     feeds,
		scheduler: scheduler,
	}
}

func (h *handler) ListGenerations(w http.ResponseWriter, r *http.Request) {
	generations, err := h.feeds.ListGenerations(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	jsonResponse(w, http.StatusOK, generations)
}

func (h *handler) ListGenerationTypes(w http.ResponseWriter, r *http.Request) {
	generationTypes, err := h.feeds.ListGenerationTypes(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	jsonResponse(w, http.StatusOK, generationTypes)
}

func (h *handler) GenerateFeed(w http.ResponseWriter, r *http.Request) {
	generationType, err := extractGenerationType(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}
	if err := h.feeds.GenerateFeed(context.Background(), generationType); err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *handler) CancelGeneration(w http.ResponseWriter, r *http.Request) {
	generationID, err := extractGenerationID(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}
	if err := h.feeds.CancelGeneration(r.Context(), generationID); err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *handler) RestartGeneration(w http.ResponseWriter, r *http.Request) {
	generationID, err := extractGenerationID(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}
	_ = generationID
	// if err := h.feeds.RestartGeneration(r.Context(), generationID); err != nil {
	// 	errorResponse(w, http.StatusInternalServerError, err)
	// 	return
	// }
	w.WriteHeader(http.StatusCreated)
}

func (h *handler) ScheduleGeneration(w http.ResponseWriter, r *http.Request) {
	generationType, err := extractGenerationType(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}
	scheduleIn := new(scheduleTaskIn)
	if err := json.NewDecoder(r.Body).Decode(scheduleIn); err != nil {
		errorResponse(w, http.StatusBadRequest, fmt.Errorf("%w: %s",
			ErrReadingRequestBody, err.Error()))
		return
	}
	if err := h.scheduleGeneration(generationType, scheduleIn); err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *handler) ListSchedules(w http.ResponseWriter, r *http.Request) {
	schedules, err := h.scheduler.ListSchedules()
	schedulesOut := makeSchedulesOut(schedules)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	jsonResponse(w, http.StatusCreated, schedulesOut)
}

func (h *handler) UnscheduleGeneration(w http.ResponseWriter, r *http.Request) {
	generationType, err := extractGenerationType(r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}
	taskID := scheduler.TaskID(generationType)
	if err := h.scheduler.RemoveTask(taskID); err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *handler) scheduleGeneration(generationType string, scheduleIn *scheduleTaskIn) error {
	taskToSchedule, err := h.makeTask(generationType, scheduleIn)
	if err != nil {
		return err
	}
	taskID := scheduler.TaskID(generationType)
	return h.scheduler.ScheduleTask(taskID, taskToSchedule)
}

func (h *handler) makeTask(generationType string, scheduleIn *scheduleTaskIn) (*scheduler.Task, error) {
	cmd, err := scheduler.NewCmd(h.feeds.GenerateFeed, context.Background(), generationType)
	if err != nil {
		return nil, err
	}
	schedule := scheduler.NewSchedule(scheduleIn.StartTimestamp, time.Second*time.Duration(scheduleIn.DelayInterval))
	taskToSchedule := scheduler.NewTask(cmd, schedule)
	return taskToSchedule, nil
}
