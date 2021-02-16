package gateway_test

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
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
			PoolSize:    5,
			ConnTimeout: time.Second,
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
				f.connection.On("Do", "").Return("", nil)
				f.connection.On("Err").Return(nil)
				f.connection.On("Do", "PING").Return("PONG", nil)
			},
		},
		{
			name:   "ping error",
			fields: defaultRedisFields(),
			setupMocks: func(f *redisFields) {
				f.dialer.On("Dial", "tcp", f.config.Addr(), mock.Anything).Return(f.connection, nil)
				f.connection.On("Do", "").Return("", nil)
				f.connection.On("Err").Return(nil)
				f.connection.On("Do", "PING").Return("", defaultErr)
			},
			wantErr: defaultErr,
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
		pool       *redis.Pool
		wantErr    error
	}{
		{
			name:   "succeed",
			fields: defaultRedisFields(),
			pool:   new(redis.Pool),
		},
		{
			name:    "redis disconnected",
			fields:  defaultRedisFields(),
			wantErr: gateway.ErrRedisDisconnected,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			redisGateway := new(gateway.RedisGateway)
			redisGateway.SetPool(testCase.pool)

			gotErr := redisGateway.Disconnect()

			assert.Equal(t, testCase.wantErr, gotErr)
		})
	}
}
