package broadcaster

import (
	"io"
	"time"

	"go-feedmaker/entity"
)

type (
	RecipientImpl   = recipient
	BroadcasterImpl = broadcaster
)

func (r *recipient) GetConn() WSConn {
	return r.conn
}

func (r *recipient) GetSend() chan []byte {
	return r.send
}

func (r *recipient) SetSend(send chan []byte) {
	r.send = send
}

func (r *recipient) GetStop() chan struct{} {
	return r.stop
}

func (r *recipient) SetStop(stop chan struct{}) {
	r.stop = stop
}

func (r *recipient) GetTicker() *time.Ticker {
	return r.ticker
}

func (r *recipient) SetTicker(ticker *time.Ticker) {
	r.ticker = ticker
}

func (r *recipient) GetOnCloseHook() CloseHook {
	return r.onCloseHook
}

func (b *broadcaster) GetRecipients() map[Recipient]bool {
	return b.recipients
}

func (b *broadcaster) SetRecipients(recipients map[Recipient]bool) {
	b.recipients = recipients
}

func (b *broadcaster) GetRegister() chan Recipient {
	return b.register
}

func (b *broadcaster) SetRegister(register chan Recipient) {
	b.register = register
}

func (b *broadcaster) GetUnregister() chan Recipient {
	return b.unregister
}

func (b *broadcaster) SetUnregister(unregister chan Recipient) {
	b.unregister = unregister
}

func (b *broadcaster) GetBroadcast() chan []byte {
	return b.broadcast
}

func (b *broadcaster) SetBroadcast(broadcast chan []byte) {
	b.broadcast = broadcast
}

func (b *broadcaster) GetStop() chan struct{} {
	return b.stop
}

func (b *broadcaster) SetStop(stop chan struct{}) {
	b.stop = stop
}

func MarshalGeneration(generation *entity.Generation, w io.Writer) error {
	return marshalGeneration(generation, w)
}
