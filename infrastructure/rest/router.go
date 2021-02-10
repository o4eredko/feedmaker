package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

type (
	Handler interface {
		ListGenerations(w http.ResponseWriter, r *http.Request)
		ListGenerationTypes(w http.ResponseWriter, r *http.Request)
		GenerateFeed(w http.ResponseWriter, r *http.Request)
		CancelGeneration(w http.ResponseWriter, r *http.Request)
	}
)

func NewRouter(handler Handler) http.Handler {
	router := mux.NewRouter()

	generations := router.PathPrefix("/generations").Subrouter()
	generations.HandleFunc("", handler.ListGenerations).Methods(http.MethodGet)
	generations.HandleFunc("/types", handler.ListGenerationTypes).Methods(http.MethodGet)
	generations.HandleFunc("/{generation-type}", handler.GenerateFeed).Methods(http.MethodPost)
	generations.HandleFunc("/{generation-id}", handler.CancelGeneration).Methods(http.MethodDelete)

	return router
}
