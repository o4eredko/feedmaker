package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"go-feedmaker/interactor"
)

type (
	handler struct {
		feeds interactor.FeedInteractor
	}
)

func NewHandler(feeds interactor.FeedInteractor) *handler {
	return &handler{
		feeds: feeds,
	}
}

func (h *handler) ListGenerations(w http.ResponseWriter, r *http.Request) {
	generations, err := h.feeds.ListGenerations(r.Context())
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	h.jsonResponse(w, http.StatusOK, generations)
}

func (h *handler) ListGenerationTypes(w http.ResponseWriter, r *http.Request) {
	generationTypes, err := h.feeds.ListGenerationTypes(r.Context())
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	h.jsonResponse(w, http.StatusOK, generationTypes)
}

func (h *handler) GenerateFeed(w http.ResponseWriter, r *http.Request) {
	generationType, err := h.extractGenerationType(r)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err)
		return
	}
	if err := h.feeds.GenerateFeed(r.Context(), generationType); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) CancelGeneration(w http.ResponseWriter, r *http.Request) {
	generationID, err := h.extractGenerationID(r)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err)
		return
	}
	if err := h.feeds.CancelGeneration(r.Context(), generationID); err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) errorResponse(w http.ResponseWriter, code int, err error) {
	body := map[string]string{"details": err.Error()}
	h.jsonResponse(w, code, body)
}

func (h *handler) jsonResponse(w http.ResponseWriter, code int, body interface{}) {
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(body); err != nil {
		log.Error().
			Err(err).
			Interface("body", body).
			Msg("send response body")
	}
}

func (h *handler) extractGenerationType(r *http.Request) (string, error) {
	return h.extractFromURL(r, "generation-type")
}

func (h *handler) extractGenerationID(r *http.Request) (string, error) {
	return h.extractFromURL(r, "generation-id")
}

func (h *handler) extractFromURL(r *http.Request, key string) (string, error) {
	vars := mux.Vars(r)
	value, found := vars[key]
	if !found {
		return "", errors.New(fmt.Sprintf("%s wasn't passed", key))
	}
	return value, nil
}
