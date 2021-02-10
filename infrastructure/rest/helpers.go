package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"go-feedmaker/infrastructure/scheduler"
	"go-feedmaker/infrastructure/scheduler/task"
)

func errorResponse(w http.ResponseWriter, code int, err error) {
	body := map[string]string{"details": err.Error()}
	jsonResponse(w, code, body)
}

func jsonResponse(w http.ResponseWriter, code int, body interface{}) {
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(body); err != nil {
		log.Error().
			Err(err).
			Interface("body", body).
			Msg("send response body")
	}
}

func extractGenerationType(r *http.Request) (string, error) {
	return extractFromURL(r, "generation-type")
}

func extractGenerationID(r *http.Request) (string, error) {
	return extractFromURL(r, "generation-id")
}

func extractFromURL(r *http.Request, key string) (string, error) {
	vars := mux.Vars(r)
	value, found := vars[key]
	if !found {
		return "", errors.New(fmt.Sprintf("%s wasn't passed", key))
	}
	return value, nil
}

func makeSchedulesOut(schedules map[scheduler.TaskID]*task.Schedule) []*scheduleOut {
	schedulesOut := make([]*scheduleOut, 0, len(schedules))
	for taskID, schedule := range schedules {
		scheduleOut := makeScheduleOut(taskID, schedule)
		schedulesOut = append(schedulesOut, scheduleOut)
	}
	return schedulesOut
}

func makeScheduleOut(taskID scheduler.TaskID, schedule *task.Schedule) *scheduleOut {
	return &scheduleOut{
		GenerationType: string(taskID),
		StartTimestamp: schedule.StartTimestamp(),
		DelayInterval:  schedule.FireInterval(),
	}
}
