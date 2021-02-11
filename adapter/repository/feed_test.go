package repository_test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

type feedFields struct {
	client *mocks.RedisClient
	conn   *mocks.Connection
	pubsub *mocks.PubSub
}

func defaultFeedFields() *feedFields {
	return &feedFields{
		client: new(mocks.RedisClient),
		conn:   new(mocks.Connection),
		pubsub: new(mocks.PubSub),
	}
}

func (f *feedFields) assertExpectations(t *testing.T) {
	f.client.AssertExpectations(t)
	f.conn.AssertExpectations(t)
	f.pubsub.AssertExpectations(t)
}

func TestFeedRepo_StoreGeneration(t *testing.T) {
	type args struct {
		ctx        context.Context
		generation *entity.Generation
	}
	testCases := []struct {
		name       string
		args       *args
		setupMocks func(*args, *feedFields)
		wantErr    error
	}{
		{
			name: "succeed",
			args: &args{
				ctx: context.Background(),
				generation: &entity.Generation{
					ID:        uuid.New().String(),
					Type:      "test",
					Progress:  13,
					StartTime: time.Now(),
					EndTime:   time.Now(),
				},
			},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				f.conn.On("Send", "MULTI").Return(nil)
				f.conn.On("Send", "SADD", mock.Anything, a.generation.ID).Return(nil)
				args := new(redis.Args).
					Add("HMSET", a.generation.ID).
					Add(mock.Anything, a.generation.Type).
					Add(mock.Anything, a.generation.Progress).
					Add(mock.Anything, a.generation.StartTime.Unix()).
					Add(mock.Anything, a.generation.EndTime.Unix())
				f.conn.On("Send", args...).Return(nil)
				f.conn.On("Do", "EXEC").Return("OK", nil)
			},
		},
		{
			name: "Do error",
			args: &args{
				ctx: context.Background(),
				generation: &entity.Generation{
					ID:        uuid.New().String(),
					Type:      "test",
					StartTime: time.Now(),
				},
			},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				f.conn.On("Send", "MULTI").Return(nil)
				f.conn.On("Send", "SADD", mock.Anything, a.generation.ID).Return(nil)
				args := new(redis.Args).
					Add("HMSET", a.generation.ID).
					Add(mock.Anything, a.generation.Type).
					Add(mock.Anything, a.generation.Progress).
					Add(mock.Anything, a.generation.StartTime.Unix())
				f.conn.On("Send", args...).Return(nil)
				f.conn.On("Do", "EXEC").Return("", defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fields := defaultFeedFields()
			tc.setupMocks(tc.args, fields)
			feedRepo := repository.NewFeedRepo(fields.client)

			gotErr := feedRepo.StoreGeneration(tc.args.ctx, tc.args.generation)

			assert.Equal(t, tc.wantErr, gotErr)
			fields.assertExpectations(t)
		})
	}
}

func TestFeedRepo_ListGenerations(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	testCases := []struct {
		name       string
		args       *args
		setupMocks func(*args, *feedFields)
		want       []*entity.Generation
		wantErr    error
	}{
		{
			name: "succeed",
			args: &args{ctx: context.Background()},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)

				f.conn.
					On("Do", "SMEMBERS", mock.Anything).
					Return([]interface{}{"123", "234"}, nil)
				f.conn.
					On("Do", "HGETALL", "123").
					Return([]interface{}{
						[]byte("type"), []byte("test1"),
						[]byte("start_time"), []byte(strconv.Itoa(int(time.Unix(1, 0).Unix()))),
					}, nil)
				f.conn.
					On("Do", "HGETALL", "234").
					Return([]interface{}{
						[]byte("type"), []byte("test2"),
						[]byte("start_time"), []byte(strconv.Itoa(int(time.Unix(11, 0).Unix()))),
						[]byte("end_time"), []byte(strconv.Itoa(int(time.Unix(20, 0).Unix()))),
					}, nil)
			},
			want: []*entity.Generation{
				{
					ID: "123", Type: "test1",
					StartTime: time.Unix(1, 0),
				},
				{
					ID: "234", Type: "test2",
					StartTime: time.Unix(11, 0),
					EndTime:   time.Unix(20, 0),
				},
			},
		},
		{
			name: "SMEMBERS error",
			args: &args{ctx: context.Background()},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				f.conn.
					On("Do", "SMEMBERS", mock.Anything).
					Return(nil, defaultErr)
			},
			wantErr: defaultErr,
		},
		{
			name: "first HGETALL error",
			args: &args{ctx: context.Background()},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				f.conn.
					On("Do", "SMEMBERS", mock.Anything).
					Return([]interface{}{"123", "234"}, nil)
				f.conn.
					On("Do", "HGETALL", "123").
					Return(nil, defaultErr)
			},
			wantErr: defaultErr,
		},
		{
			name: "second HGETALL error",
			args: &args{ctx: context.Background()},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				f.conn.
					On("Do", "SMEMBERS", mock.Anything).
					Return([]interface{}{"123", "234"}, nil)
				f.conn.
					On("Do", "HGETALL", "123").
					Return([]interface{}{
						[]byte("type"), []byte("test1"),
						[]byte("start_time"), []byte(strconv.Itoa(int(time.Unix(1, 0).Unix()))),
						[]byte("end_time"), []byte(strconv.Itoa(int(time.Unix(10, 0).Unix()))),
					}, nil)
				f.conn.
					On("Do", "HGETALL", "234").Return(nil, defaultErr)
			},
			wantErr: defaultErr,
		},
		{
			name: "invalid timestamp error",
			args: &args{ctx: context.Background()},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				f.conn.
					On("Do", "SMEMBERS", mock.Anything).
					Return([]interface{}{"123"}, nil)
				f.conn.
					On("Do", "HGETALL", "123").
					Return([]interface{}{
						[]byte("type"), []byte("test1"),
						[]byte("start_time"), []byte("invalid"),
					}, nil)
			},
			wantErr: entity.ErrInvalidTimestamp,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fields := defaultFeedFields()
			tc.setupMocks(tc.args, fields)
			feedRepo := repository.NewFeedRepo(fields.client)

			got, gotErr := feedRepo.ListGenerations(tc.args.ctx)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, tc.want, got)
			fields.assertExpectations(t)
		})
	}
}

func TestFeedRepo_UpdateProgress(t *testing.T) {
	type args struct {
		ctx        context.Context
		generation *entity.Generation
	}
	testCases := []struct {
		name       string
		args       *args
		setupMocks func(*args, *feedFields)
		wantErr    error
	}{
		{
			name: "succeed",
			args: &args{
				ctx: context.Background(),
				generation: &entity.Generation{
					ID:        uuid.New().String(),
					Type:      "test",
					Progress:  100,
					StartTime: time.Now(),
					EndTime:   time.Now(),
				},
			},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				args := new(redis.Args).
					Add("HSET", a.generation.ID).
					Add(mock.Anything, a.generation.Progress).
					Add(mock.Anything, a.generation.EndTime.Unix())
				f.conn.On("Do", args...).Return("", nil)

				args = new(redis.Args).Add("PUBLISH", a.generation.ID, a.generation.Progress)
				f.conn.On("Do", args...).Return("", nil)
			},
		},
		{
			name: "succeed without end time",
			args: &args{
				ctx: context.Background(),
				generation: &entity.Generation{
					ID:        uuid.New().String(),
					Type:      "test",
					Progress:  80,
					StartTime: time.Now(),
				},
			},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				args := new(redis.Args).Add("HSET", a.generation.ID).Add(mock.Anything, a.generation.Progress)
				f.conn.On("Do", args...).Return("", nil)

				args = new(redis.Args).Add("PUBLISH", a.generation.ID, a.generation.Progress)
				f.conn.On("Do", args...).Return("", nil)
			},
		},
		{
			name: "HSET error",
			args: &args{
				ctx: context.Background(),
				generation: &entity.Generation{
					ID:        uuid.New().String(),
					Type:      "test",
					Progress:  100,
					StartTime: time.Now(),
					EndTime:   time.Now(),
				},
			},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				args := new(redis.Args).
					Add("HSET", a.generation.ID).
					Add(mock.Anything, a.generation.Progress).
					Add(mock.Anything, a.generation.EndTime.Unix())
				f.conn.On("Do", args...).Return("", defaultErr)
			},
			wantErr: defaultErr,
		},
		{
			name: "PUBLISH error",
			args: &args{
				ctx: context.Background(),
				generation: &entity.Generation{
					ID:        uuid.New().String(),
					Type:      "test",
					Progress:  100,
					StartTime: time.Now(),
					EndTime:   time.Now(),
				},
			},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				args := new(redis.Args).
					Add("HSET", a.generation.ID).
					Add(mock.Anything, a.generation.Progress).
					Add(mock.Anything, a.generation.EndTime.Unix())
				f.conn.On("Do", args...).Return("", nil)

				args = new(redis.Args).Add("PUBLISH", a.generation.ID, a.generation.Progress)
				f.conn.On("Do", args...).Return("", defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fields := defaultFeedFields()
			tc.setupMocks(tc.args, fields)
			feedRepo := repository.NewFeedRepo(fields.client)

			gotErr := feedRepo.UpdateProgress(tc.args.ctx, tc.args.generation)

			assert.Equal(t, tc.wantErr, gotErr)
			fields.assertExpectations(t)
		})
	}
}

func TestFeedRepo_CancelGeneration(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	testCases := []struct {
		name       string
		args       *args
		setupMocks func(*args, *feedFields)
		wantErr    error
	}{
		{
			name: "succeed",
			args: &args{
				ctx: context.Background(),
				id:  uuid.NewString(),
			},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				f.conn.
					On("Do", "PUBLISH", fmt.Sprintf("%s.canceled", a.id), mock.Anything).
					Return("", nil)
			},
		},
		{
			name: "Do error",
			args: &args{
				ctx: context.Background(),
				id:  uuid.NewString(),
			},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				f.conn.
					On("Do", "PUBLISH", fmt.Sprintf("%s.canceled", a.id), mock.Anything).
					Return("", defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fields := defaultFeedFields()
			tc.setupMocks(tc.args, fields)
			feedRepo := repository.NewFeedRepo(fields.client)

			gotErr := feedRepo.CancelGeneration(tc.args.ctx, tc.args.id)

			assert.Equal(t, tc.wantErr, gotErr)
			fields.assertExpectations(t)
		})
	}
}