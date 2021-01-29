package gateway_test

import (
	"testing"
	"time"

	"github.com/mediocregopher/radix/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-feedmaker/infrastructure/gateway"
	"go-feedmaker/infrastructure/gateway/mocks"
)

type redisFields struct {
	dialer     *mocks.RedisDialer
	config     gateway.RedisConfig
	connection *mocks.RedisConnection
}

func defaultRedisFields() *redisFields {
	return &redisFields{
		dialer: new(mocks.RedisDialer),
		config: gateway.RedisConfig{
			Host:        "localhost",
			Port:        "5000",
			ConnTimeout: time.Second,
			PoolSize:    2,
		},
		connection: new(mocks.RedisConnection),
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

func TestRedisGateway_Do(t *testing.T) {
	type args struct {
		action radix.Action
	}
	defaultArgs := func() *args {
		return &args{radix.Cmd(nil, "test", "test")}
	}
	testCases := []struct {
		name       string
		args       *args
		fields     *redisFields
		setupMocks func(*args, *redisFields)
		wantErr    error
	}{
		{
			name:   "succeed",
			args:   defaultArgs(),
			fields: defaultRedisFields(),
			setupMocks: func(a *args, f *redisFields) {
				f.connection.On("Do", a.action).Return(nil)
			},
		},
		{
			name: "disconnected error",
			args: defaultArgs(),
			fields: &redisFields{
				dialer:     new(mocks.RedisDialer),
				config:     gateway.RedisConfig{},
				connection: nil,
			},
			setupMocks: func(a *args, f *redisFields) {},
			wantErr:    gateway.ErrRedisDisconnected,
		},
		{
			name:   "Do error",
			args:   defaultArgs(),
			fields: defaultRedisFields(),
			setupMocks: func(a *args, f *redisFields) {
				f.connection.On("Do", a.action).Return(defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.args, testCase.fields)
			redisGateway := gateway.RedisGateway{
				Dialer: testCase.fields.dialer,
				Config: testCase.fields.config,
			}
			redisGateway.SetConnection(testCase.fields.connection)

			gotErr := redisGateway.Do(testCase.args.action)

			assert.Equal(t, testCase.wantErr, gotErr)
			testCase.fields.dialer.AssertExpectations(t)
			if testCase.fields.connection != nil {
				testCase.fields.connection.AssertExpectations(t)
			}
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
