package gateway_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-feedmaker/infrastructure/gateway"
	"go-feedmaker/infrastructure/gateway/mocks"
)

type redisFields struct {
	dialer     *mocks.RedisDialer
	config     gateway.RedisConfig
	connection *mocks.Connection
}

func defaultRedisFields() *redisFields {
	return &redisFields{
		dialer: new(mocks.RedisDialer),
		config: gateway.RedisConfig{
			Host:        "localhost",
			Port:        "5000",
			ConnTimeout: time.Second,
		},
		connection: new(mocks.Connection),
	}
}

func TestRedisGateway_Connect(t *testing.T) {
	testCases := []struct {
		name       string
		fields     *redisFields
		setupMocks func(f *redisFields)
		wantErr    error
	}{
		{
			name:   "succeed",
			fields: defaultRedisFields(),
			setupMocks: func(f *redisFields) {
				f.dialer.On("Dial", "tcp", f.config.Addr(), mock.Anything).Return(f.connection, nil)
			},
		},
		{
			name:   "dial error",
			fields: defaultRedisFields(),
			setupMocks: func(f *redisFields) {
				f.dialer.On("Dial", "tcp", f.config.Addr(), mock.Anything).Return(nil, defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields)
			redisGateway := gateway.RedisGateway{
				Dialer: testCase.fields.dialer,
				Config: testCase.fields.config,
			}

			gotErr := redisGateway.Connect()

			assert.Equal(t, testCase.wantErr, gotErr)
			testCase.fields.dialer.AssertExpectations(t)
			testCase.fields.connection.AssertExpectations(t)
		})
	}
}

func TestRedisGateway_Disconnect(t *testing.T) {
	testCases := []struct {
		name       string
		fields     *redisFields
		setupMocks func(f *redisFields)
		wantErr    error
	}{
		{
			name:   "succeed",
			fields: defaultRedisFields(),
			setupMocks: func(f *redisFields) {
				f.connection.On("Close").Return(nil)
			},
		},
		{
			name: "disconnected error",
			fields: &redisFields{
				dialer:     new(mocks.RedisDialer),
				config:     gateway.RedisConfig{},
				connection: nil,
			},
			setupMocks: func(f *redisFields) {},
			wantErr:    gateway.ErrRedisDisconnected,
		},
		{
			name:   "Close error",
			fields: defaultRedisFields(),
			setupMocks: func(f *redisFields) {
				f.connection.On("Close").Return(defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields)
			redisGateway := gateway.RedisGateway{
				Dialer: testCase.fields.dialer,
				Config: testCase.fields.config,
			}
			redisGateway.SetConnection(testCase.fields.connection)

			gotErr := redisGateway.Disconnect()

			assert.Equal(t, testCase.wantErr, gotErr)
			testCase.fields.dialer.AssertExpectations(t)
			if testCase.fields.connection != nil {
				testCase.fields.connection.AssertExpectations(t)
			}
		})
	}
}
