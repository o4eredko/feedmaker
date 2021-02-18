package rest

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type (
	Handler interface {
		ListGenerations(w http.ResponseWriter, r *http.Request)
		ListGenerationTypes(w http.ResponseWriter, r *http.Request)
		GenerateFeed(w http.ResponseWriter, r *http.Request)
		CancelGeneration(w http.ResponseWriter, r *http.Request)
		RestartGeneration(w http.ResponseWriter, r *http.Request)
		ScheduleGeneration(w http.ResponseWriter, r *http.Request)
		ListSchedules(w http.ResponseWriter, r *http.Request)
		UnscheduleGeneration(w http.ResponseWriter, r *http.Request)
	}

	WSHandler interface {
		ServeWS(w http.ResponseWriter, r *http.Request)
	}
)

func NewRouter(handler Handler, wsHandler WSHandler) http.Handler {
	router := mux.NewRouter()

	generations := router.PathPrefix("/generations").Subrouter()

	generations.HandleFunc("", handler.ListGenerations).Methods(http.MethodGet)

	generations.HandleFunc("/types", handler.ListGenerationTypes).Methods(http.MethodGet)

	generations.HandleFunc("/types/{generation-type}", handler.GenerateFeed).Methods(http.MethodPost)

	generations.HandleFunc("/id/{generation-id}", handler.RestartGeneration).Methods(http.MethodPost)
	generations.HandleFunc("/id/{generation-id}", handler.CancelGeneration).Methods(http.MethodDelete)

	generations.HandleFunc("/schedules", handler.ListSchedules).Methods(http.MethodGet)
	generations.HandleFunc("/types/{generation-type}/schedules", handler.ScheduleGeneration).Methods(http.MethodPost)
	generations.HandleFunc("/types/{generation-type}/schedules", handler.UnscheduleGeneration).Methods(http.MethodDelete)

	ws := router.PathPrefix("/ws").Subrouter()
	ws.HandleFunc("/progress", wsHandler.ServeWS)

	headersOK := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Authorization"})
	originsOK := handlers.AllowedOrigins([]string{"*"})
	methodsOK := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"})

	return handlers.CORS(headersOK, originsOK, methodsOK)(router)
}
