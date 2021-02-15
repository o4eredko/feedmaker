package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
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
