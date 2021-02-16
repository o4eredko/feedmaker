package broadcaster_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-feedmaker/infrastructure/rest/broadcaster"
	"go-feedmaker/infrastructure/rest/mocks"
)

type (
	broadcasterFields struct {
		recipients map[broadcaster.Recipient]bool
		register   chan broadcaster.Recipient
		unregister chan broadcaster.Recipient
		broadcast  chan []byte
		stop       chan struct{}
	}

	mockedRecipients []*mocks.Recipient
)

const (
	recipientsAmount   = 5
	recipientToActions = 3
)

func (m mockedRecipients) AssertExpectations(t *testing.T) {
	for _, recipient := range m {
		recipient.AssertExpectations(t)
	}
}

func defaultBroadcasterFields() *broadcasterFields {
	return &broadcasterFields{
		recipients: make(map[broadcaster.Recipient]bool),
		register:   make(chan broadcaster.Recipient),
		unregister: make(chan broadcaster.Recipient),
		broadcast:  make(chan []byte),
		stop:       make(chan struct{}),
	}
}

func makeMockedRecipients(amount int) mockedRecipients {
	recipients := make([]*mocks.Recipient, amount)
	for i := 0; i < amount; i++ {
		recipients[i] = new(mocks.Recipient)
	}
	return recipients
}

func TestNewBroadcaster(t *testing.T) {
	b := broadcaster.NewBroadcaster()
	assert.NotNil(t, b.GetRecipients())
	assert.NotNil(t, b.GetRegister())
	assert.NotNil(t, b.GetUnregister())
	assert.NotNil(t, b.GetBroadcast())
	assert.NotNil(t, b.GetStop())
}

func Test_broadcaster_Start(t *testing.T) {
	testCases := []struct {
		name       string
		fields     *broadcasterFields
		setupMocks func([]*mocks.Recipient)
		do         func(*broadcaster.BroadcasterImpl, *broadcasterFields, []*mocks.Recipient)
	}{
		{
			name:   "register",
			fields: defaultBroadcasterFields(),
			setupMocks: func(recipients []*mocks.Recipient) {
				recipient := recipients[recipientToActions]
				recipient.On("Start")
				recipient.On("OnCloseHook", mock.Anything)
			},
			do: func(b *broadcaster.BroadcasterImpl, fields *broadcasterFields, recipients []*mocks.Recipient) {
				b.SetRegister(fields.register)
				fields.register <- recipients[recipientToActions]
				<-time.After(time.Millisecond)
			},
		},
		{
			name:   "unregister",
			fields: defaultBroadcasterFields(),
			setupMocks: func(recipients []*mocks.Recipient) {
				recipient := recipients[recipientToActions]
				recipient.On("Stop")
			},
			do: func(b *broadcaster.BroadcasterImpl, fields *broadcasterFields, recipients []*mocks.Recipient) {
				b.SetUnregister(fields.unregister)
				fields.unregister <- recipients[recipientToActions]
				<-time.After(time.Millisecond)
			},
		},
		{
			name:   "broadcast",
			fields: defaultBroadcasterFields(),
			setupMocks: func(recipients []*mocks.Recipient) {
				for _, recipient := range recipients {
					recipient.On("Send", defaultMessage)
				}
			},
			do: func(b *broadcaster.BroadcasterImpl, fields *broadcasterFields, recipients []*mocks.Recipient) {
				recipientsMap := make(map[broadcaster.Recipient]bool)
				for _, recipient := range recipients {
					recipientsMap[recipient] = true
				}
				b.SetRecipients(recipientsMap)
				b.SetBroadcast(fields.broadcast)
				fields.broadcast <- defaultMessage
				<-time.After(time.Millisecond)
			},
		},
		{
			name:   "stop",
			fields: defaultBroadcasterFields(),
			setupMocks: func(recipients []*mocks.Recipient) {
				for _, recipient := range recipients {
					recipient.On("Stop")
				}
			},
			do: func(b *broadcaster.BroadcasterImpl, fields *broadcasterFields, recipients []*mocks.Recipient) {
				recipientsMap := make(map[broadcaster.Recipient]bool)
				for _, recipient := range recipients {
					recipientsMap[recipient] = true
				}
				b.SetRecipients(recipientsMap)
				b.SetStop(fields.stop)
				fields.stop <- struct{}{}
				<-time.After(time.Millisecond)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			recipients := makeMockedRecipients(recipientsAmount)
			testCase.setupMocks(recipients)
			b := broadcaster.NewBroadcaster()
			go b.Start()
			testCase.do(b, testCase.fields, recipients)
			recipients.AssertExpectations(t)
		})
	}
}

func Test_broadcaster_Stop(t *testing.T) {
	fields := defaultBroadcasterFields()
	b := broadcaster.NewBroadcaster()
	b.SetStop(fields.stop)
	go b.Stop()
	gotEvent := <-fields.stop
	assert.Equal(t, struct{}{}, gotEvent)
}

func Test_broadcaster_Register(t *testing.T) {
	fields := defaultBroadcasterFields()
	b := broadcaster.NewBroadcaster()
	b.SetRegister(fields.register)
	recipient := new(mocks.Recipient)
	go b.Register(recipient)
	gotRecipient := <-fields.register
	assert.Equal(t, recipient, gotRecipient)
}

func Test_broadcaster_Unregister(t *testing.T) {
	fields := defaultBroadcasterFields()
	b := broadcaster.NewBroadcaster()
	b.SetUnregister(fields.unregister)
	recipient := new(mocks.Recipient)
	go b.Unregister(recipient)
	gotRecipient := <-fields.unregister
	assert.Equal(t, recipient, gotRecipient)
}

func Test_broadcaster_Broadcast(t *testing.T) {
	fields := defaultBroadcasterFields()
	b := broadcaster.NewBroadcaster()
	b.SetBroadcast(fields.broadcast)
	go b.Broadcast(defaultMessage)
	gotMessage := <-fields.broadcast
	assert.Equal(t, defaultMessage, gotMessage)
}
