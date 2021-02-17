package broadcaster

type (
	broadcaster struct {
		recipients map[Recipient]bool
		register   chan Recipient
		unregister chan Recipient
		broadcast  chan []byte
		stop       chan struct{}
	}

	Recipient interface {
		Start()
		Stop()
		Send([]byte)
		OnCloseHook(hook CloseHook)
	}
)

func NewBroadcaster() *broadcaster {
	return &broadcaster{
		recipients: make(map[Recipient]bool),
		register:   make(chan Recipient),
		unregister: make(chan Recipient),
		broadcast:  make(chan []byte),
		stop:       make(chan struct{}),
	}
}

func (b *broadcaster) Start() {
	for {
		select {
		case recipient := <-b.register:
			b.pushRecipient(recipient)
		case recipient := <-b.unregister:
			b.stopRecipient(recipient)
		case msg := <-b.broadcast:
			b.broadcastMsg(msg)
		case <-b.stop:
			b.stopBroadcasting()
			return
		}
	}
}

func (b *broadcaster) Stop() {
	b.stop <- struct{}{}
}

func (b *broadcaster) Register(recipient Recipient) {
	b.register <- recipient
}

func (b *broadcaster) Unregister(recipient Recipient) {
	b.unregister <- recipient
}

func (b *broadcaster) Broadcast(msg []byte) {
	b.broadcast <- msg
}

func (b *broadcaster) pushRecipient(recipient Recipient) {
	b.recipients[recipient] = true
	hook := b.makeOnCloseHook(recipient)
	recipient.OnCloseHook(hook)
	go recipient.Start()
}

func (b *broadcaster) stopRecipient(recipient Recipient) {
	recipient.Stop()
}

func (b *broadcaster) broadcastMsg(msg []byte) {
	for recipient := range b.recipients {
		recipient.Send(msg)
	}
}

func (b *broadcaster) stopBroadcasting() {
	close(b.register)
	close(b.broadcast)
	for recipient := range b.recipients {
		recipient.Stop()
	}
	close(b.unregister)
}

func (b *broadcaster) makeOnCloseHook(recipient Recipient) CloseHook {
	return func() {
		if _, ok := b.recipients[recipient]; ok {
			delete(b.recipients, recipient)
		}
	}
}
