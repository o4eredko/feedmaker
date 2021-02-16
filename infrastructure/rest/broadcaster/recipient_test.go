package broadcaster_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"

	"go-feedmaker/infrastructure/rest/broadcaster"
	"go-feedmaker/infrastructure/rest/mocks"
)

type (
	recipientFields struct {
		conn      *mocks.WSConn
		send      chan []byte
		stop      chan struct{}
		ticker    *time.Ticker
		closeHook *mocks.CloseHook
	}
)

var (
	defaultMessage = []byte("default message")
)

func (f *recipientFields) AssertExpectations(t *testing.T) {
	f.conn.AssertExpectations(t)
	f.closeHook.AssertExpectations(t)
}

func defaultRecipientFields() *recipientFields {
	return &recipientFields{
		conn:      new(mocks.WSConn),
		send:      make(chan []byte),
		stop:      make(chan struct{}),
		ticker:    &time.Ticker{C: make(chan time.Time)},
		closeHook: new(mocks.CloseHook),
	}
}

func TestNewRecipient(t *testing.T) {
	fields := defaultRecipientFields()
	recipient := broadcaster.NewRecipient(fields.conn)
	assert.Equal(t, fields.conn, recipient.GetConn())
	assert.NotNil(t, recipient.GetSend())
	assert.NotNil(t, recipient.GetStop())
	assert.NotNil(t, recipient.GetTicker())
	assert.Nil(t, recipient.GetOnCloseHook())
}

func Test_recipient_OnCloseHook(t *testing.T) {
	fields := defaultRecipientFields()
	recipient := broadcaster.NewRecipient(fields.conn)
	recipient.OnCloseHook(fields.closeHook.Execute)
	wantFunc := reflect.ValueOf(fields.closeHook.Execute)
	gotFunc := reflect.ValueOf(recipient.GetOnCloseHook())
	assert.Equal(t, wantFunc.Pointer(), gotFunc.Pointer(),
		"want: %v\ngot: %v", wantFunc.String(), gotFunc.String())
}

func Test_recipient_Start(t *testing.T) {
	testCases := []struct {
		name       string
		fields     *recipientFields
		setupMocks func(*recipientFields)
		do         func(*broadcaster.RecipientImpl, *recipientFields)
	}{
		{
			name:   "send message",
			fields: defaultRecipientFields(),
			setupMocks: func(fields *recipientFields) {
				fields.conn.On("WriteMessage", websocket.TextMessage, defaultMessage).Return(nil)
			},
			do: func(recipient *broadcaster.RecipientImpl, fields *recipientFields) {
				recipient.SetSend(fields.send)
				fields.send <- defaultMessage
			},
		},
		{
			name:   "ping",
			fields: defaultRecipientFields(),
			setupMocks: func(fields *recipientFields) {
				fields.conn.On("WriteMessage", websocket.PingMessage, []byte(nil)).Return(nil)
			},
			do: func(recipient *broadcaster.RecipientImpl, fields *recipientFields) {
				tickerChan := make(chan time.Time)
				ticker := &time.Ticker{C: tickerChan}
				recipient.SetTicker(ticker)
				tickerChan <- time.Now().UTC()
			},
		},
		{
			name:   "stop",
			fields: defaultRecipientFields(),
			setupMocks: func(fields *recipientFields) {
				fields.conn.On("Close").Return(nil)
			},
			do: func(recipient *broadcaster.RecipientImpl, fields *recipientFields) {
				recipient.Stop()
				<-time.After(time.Millisecond)
			},
		},
		{
			name:   "stop with hook",
			fields: defaultRecipientFields(),
			setupMocks: func(fields *recipientFields) {
				fields.closeHook.On("Execute")
				fields.conn.On("Close").Return(nil)
			},
			do: func(recipient *broadcaster.RecipientImpl, fields *recipientFields) {
				recipient.OnCloseHook(fields.closeHook.Execute)
				recipient.Stop()
				<-time.After(time.Millisecond)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields)
			r := broadcaster.NewRecipient(testCase.fields.conn)
			go r.Start()
			testCase.do(r, testCase.fields)
			testCase.fields.AssertExpectations(t)
		})
	}
}

func Test_recipient_Send(t *testing.T) {
	fields := defaultRecipientFields()
	r := broadcaster.NewRecipient(fields.conn)
	r.SetSend(fields.send)
	go r.Send(defaultMessage)
	gotMsg := <-fields.send
	assert.Equal(t, defaultMessage, gotMsg)
	fields.AssertExpectations(t)
}

func Test_recipient_Stop(t *testing.T) {
	fields := defaultRecipientFields()
	r := broadcaster.NewRecipient(fields.conn)
	r.SetStop(fields.stop)
	go r.Stop()
	gotEvent := <-fields.stop
	assert.Equal(t, struct{}{}, gotEvent)
	fields.AssertExpectations(t)
}
