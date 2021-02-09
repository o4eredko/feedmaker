package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-feedmaker/adapter/repository"
	"go-feedmaker/adapter/repository/mocks"
	"go-feedmaker/entity"
)

var (
	defaultErr = errors.New("test error")
)

func TestFeedRepo_StoreGeneration(t *testing.T) {
	type args struct {
		ctx        context.Context
		generation *entity.Generation
	}
	testCases := []struct {
		name       string
		args       *args
		setupMocks func(*args, *mocks.RedisClient)
		wantErr    error
	}{
		{
			name: "succeed",
			args: &args{
				ctx: context.Background(),
				generation: &entity.Generation{
					ID:        uuid.New().String(),
					Type:      "test",
					StartTime: time.Now(),
				},
			},
			setupMocks: func(a *args, client *mocks.RedisClient) {
				client.On("Send", "MULTI").Return(nil)
				client.On("Send", "SADD", mock.Anything, a.generation.ID).Return(nil)
				args := new(redis.Args).Add("HMSET").Add(a.generation.ID).AddFlat(a.generation)
				client.On("Send", args...).Return(nil)
				client.On("Do", "EXEC").Return("OK", nil)
			},
		},
		{
			name: "client.Do error",
			args: &args{
				ctx: context.Background(),
				generation: &entity.Generation{
					ID:        uuid.New().String(),
					Type:      "test",
					StartTime: time.Now(),
				},
			},
			setupMocks: func(a *args, client *mocks.RedisClient) {
				client.On("Send", "MULTI").Return(nil)
				client.On("Send", "SADD", mock.Anything, a.generation.ID).Return(nil)
				args := new(redis.Args).Add("HMSET").Add(a.generation.ID).AddFlat(a.generation)
				client.On("Send", args...).Return(nil)
				client.On("Do", "EXEC").Return("", defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := new(mocks.RedisClient)
			tc.setupMocks(tc.args, client)
			feedRepo := repository.NewFeedRepo(client)

			gotErr := feedRepo.StoreGeneration(tc.args.ctx, tc.args.generation)

			assert.Equal(t, tc.wantErr, gotErr)
			client.AssertExpectations(t)
		})
	}
}
