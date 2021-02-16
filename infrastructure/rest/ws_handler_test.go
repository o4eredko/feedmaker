package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-feedmaker/infrastructure/rest"
	"go-feedmaker/infrastructure/rest/mocks"
)

type (
	wsHandlerFields struct {
		upgrader    *mocks.Upgrader
		broadcaster *mocks.Broadcaster
	}
)

var (
	defaultConn = new(websocket.Conn)
)

func defaultWSHandlerFields() *wsHandlerFields {
	return &wsHandlerFields{
		upgrader:    new(mocks.Upgrader),
		broadcaster: new(mocks.Broadcaster),
	}
}

func (f *wsHandlerFields) AssertExpectations(t *testing.T) {
	f.upgrader.AssertExpectations(t)
	f.broadcaster.AssertExpectations(t)
}

func TestNewWSHandler(t *testing.T) {
	fields := defaultWSHandlerFields()
	h := rest.NewWSHandler(fields.upgrader, fields.broadcaster)
	assert.Equal(t, fields.upgrader, h.Upgrader())
	assert.Equal(t, fields.broadcaster, h.Broadcaster())
}

func Test_wsHandler_ServeWS(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	defaultArgs := func() *args {
		return &args{
			w: httptest.NewRecorder(),
			r: &http.Request{},
		}
	}
	testCases := []struct {
		name           string
		fields         *wsHandlerFields
		args           *args
		setupMocks     func(*wsHandlerFields, *args)
		wantStatusCode int
	}{
		{
			name:   "succeed",
			fields: defaultWSHandlerFields(),
			args:   defaultArgs(),
			setupMocks: func(fields *wsHandlerFields, args *args) {
				fields.upgrader.On("Upgrade", mock.Anything, mock.Anything, mock.Anything).Return(defaultConn, nil)
				fields.broadcaster.On("Register", mock.Anything)
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:   "upgrade error",
			fields: defaultWSHandlerFields(),
			args:   defaultArgs(),
			setupMocks: func(fields *wsHandlerFields, args *args) {
				fields.upgrader.On("Upgrade", mock.Anything, mock.Anything, mock.Anything).Return(nil, defaultTestErr)
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields, testCase.args)
			h := rest.NewWSHandler(testCase.fields.upgrader, testCase.fields.broadcaster)
			h.ServeWS(testCase.args.w, testCase.args.r)
			gotStatusCode := testCase.args.w.Code
			assert.Equal(t, testCase.wantStatusCode, gotStatusCode)
			testCase.fields.AssertExpectations(t)
		})
	}
}
