package rest

import (
	"net/http"

	"github.com/gorilla/websocket"

	"go-feedmaker/infrastructure/rest/broadcaster"
)

type (
	wsHandler struct {
		upgrader    Upgrader
		broadcaster Broadcaster
	}

	Upgrader interface {
		Upgrade(http.ResponseWriter, *http.Request, http.Header) (*websocket.Conn, error)
	}

	Broadcaster interface {
		Register(recipient broadcaster.Recipient)
	}
)

func NewWSHandler(upgrader Upgrader, broadcaster Broadcaster) *wsHandler {
	return &wsHandler{
		upgrader:    upgrader,
		broadcaster: broadcaster,
	}
}

func (h *wsHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}
	recipient := broadcaster.NewRecipient(conn)
	h.broadcaster.Register(recipient)
}
