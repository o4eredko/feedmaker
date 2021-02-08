package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type (
	APIServer struct {
		server *http.Server
	}
)

func NewAPIServer(config *Config, router http.Handler) *APIServer {
	server := &http.Server{
		Addr:    config.Addr(),
		Handler: router,
	}
	return &APIServer{
		server: server,
	}
}

func (s *APIServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *APIServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("APIServer.Stop()")
	}
}
