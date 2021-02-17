package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"go-feedmaker/infrastructure/scheduler"
)

var (
	ErrValueNotFoundInURL = errors.New("not found in url")
	ErrReadingRequestBody = errors.New("reading request body")
)

func errorResponse(w http.ResponseWriter, code int, err error) {
	body := map[string]string{"details": err.Error()}
	jsonResponse(w, code, body)
}

func jsonResponse(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
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
		return "", fmt.Errorf("looking for %v: %w", key, ErrValueNotFoundInURL)
	}
	return value, nil
}

func makeSchedulesOut(schedules map[scheduler.TaskID]*scheduler.Schedule) map[scheduler.TaskID]*scheduleOut {
	schedulesOut := make(map[scheduler.TaskID]*scheduleOut, len(schedules))
	for taskID, schedule := range schedules {
		scheduleOut := makeScheduleOut(schedule)
		schedulesOut[taskID] = scheduleOut
	}
	return schedulesOut
}

func makeScheduleOut(schedule *scheduler.Schedule) *scheduleOut {
	return &scheduleOut{
		StartTimestamp: schedule.StartTimestamp().Format(time.RFC3339),
		DelayInterval:  int(schedule.FireInterval().Seconds()),
	}
}
