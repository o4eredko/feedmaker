package broadcaster

import (
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type (
	recipient struct {
		conn        WSConn
		send        chan []byte
		stop        chan struct{}
		ticker      *time.Ticker
		onCloseHook CloseHook
	}

	WSConn interface {
		RemoteAddr() net.Addr
		WriteMessage(int, []byte) error
		Close() error
	}

	CloseHook func()
)

const (
	tickInterval = time.Second
)

func NewRecipient(conn WSConn) *recipient {
	return &recipient{
		conn:   conn,
		send:   make(chan []byte),
		stop:   make(chan struct{}),
		ticker: time.NewTicker(tickInterval),
	}
}

func (r *recipient) OnCloseHook(hook CloseHook) {
	r.onCloseHook = hook
}

func (r *recipient) Start() {
	for {
		select {
		case msg := <-r.send:
			r.sendMsg(msg)
		case <-r.ticker.C:
			r.ping()
		case <-r.stop:
			r.stopSending()
			return
		}
	}
}

func (r *recipient) Send(msg []byte) {
	r.send <- msg
}

func (r *recipient) Stop() {
	r.stop <- struct{}{}
}

func (r *recipient) sendMsg(msg []byte) {
	if err := r.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Error().
			Err(err).
			Str("remote_addr", r.conn.RemoteAddr().String()).
			Str("message", string(msg)).
			Msg("send message")
	}
}

func (r *recipient) ping() {
	if err := r.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		r.Stop()
		log.Error().
			Err(err).
			Str("remote_addr", r.conn.RemoteAddr().String()).
			Msg("ping")
	}
}

func (r *recipient) stopSending() {
	close(r.send)
	r.ticker.Stop()
	close(r.stop)
	if r.onCloseHook != nil {
		r.onCloseHook()
	}
	if err := r.conn.Close(); err != nil {
		log.Error().
			Err(err).
			Str("remote_addr", r.conn.RemoteAddr().String()).
			Msg("close connection")
	}
}
