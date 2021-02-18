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
	helper "go-feedmaker/infrastructure/testing"
)

var (
	defaultErr = errors.New("test error")
)

type feedFields struct {
	config map[string]*repository.FeedConfig
	client *mocks.RedisClient
	conn   *mocks.Connection
	pubsub *mocks.PubSub
	sql    *mocks.SqlGateway
	ftp    *mocks.FtpGateway
}

func defaultFeedFields() *feedFields {
	return &feedFields{
		config: make(map[string]*repository.FeedConfig, 1),
		client: new(mocks.RedisClient),
		conn:   new(mocks.Connection),
		pubsub: new(mocks.PubSub),
		sql:    new(mocks.SqlGateway),
		ftp:    new(mocks.FtpGateway),
	}
}

func (f *feedFields) assertExpectations(t *testing.T) {
	assert.True(t, f.client.AssertExpectations(t))
	assert.True(t, f.conn.AssertExpectations(t))
	assert.True(t, f.pubsub.AssertExpectations(t))
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
					Add(mock.Anything, a.generation.DataFetched).
					Add(mock.Anything, a.generation.FilesUploaded).
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
					Add(mock.Anything, a.generation.DataFetched).
					Add(mock.Anything, a.generation.FilesUploaded).
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
			feedRepo := repository.NewFeedRepo(fields.config, fields.client, fields.ftp)

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
						[]byte("progress"), []byte("100"),
						[]byte("files_uploaded"), []byte("4"),
						[]byte("data_fetched"), []byte("1"),
						[]byte("is_canceled"), []byte("1"),
						[]byte("start_time"), []byte(strconv.Itoa(int(time.Unix(1, 0).Unix()))),
					}, nil)
				f.conn.
					On("Do", "HGETALL", "234").
					Return([]interface{}{
						[]byte("type"), []byte("test2"),
						[]byte("progress"), []byte("43"),
						[]byte("files_uploaded"), []byte("5"),
						[]byte("data_fetched"), []byte("0"),
						[]byte("start_time"), []byte(strconv.Itoa(int(time.Unix(11, 0).Unix()))),
						[]byte("end_time"), []byte(strconv.Itoa(int(time.Unix(20, 0).Unix()))),
					}, nil)
			},
			want: []*entity.Generation{
				{
					ID: "123", Type: "test1",
					Progress:      100,
					FilesUploaded: 4,
					DataFetched:   true,
					IsCanceled:    true,
					StartTime:     time.Unix(1, 0),
				},
				{
					ID: "234", Type: "test2",
					Progress:      43,
					FilesUploaded: 5,
					DataFetched:   false,
					StartTime:     time.Unix(11, 0),
					EndTime:       time.Unix(20, 0),
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
			feedRepo := repository.NewFeedRepo(fields.config, fields.client, fields.ftp)

			got, gotErr := feedRepo.ListGenerations(tc.args.ctx)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, tc.want, got)
			fields.assertExpectations(t)
		})
	}
}

func TestFeedRepo_UpdateGenerationState(t *testing.T) {
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
					Add(mock.Anything, a.generation.DataFetched).
					Add(mock.Anything, a.generation.IsCanceled).
					Add(mock.Anything, a.generation.FilesUploaded).
					Add(mock.Anything, a.generation.EndTime.Unix())
				f.conn.On("Do", args...).Return("", nil)

				args = new(redis.Args).Add("PUBLISH", "generation.updated", a.generation.ID)
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

				args := new(redis.Args).
					Add("HSET", a.generation.ID).
					Add(mock.Anything, a.generation.Progress).
					Add(mock.Anything, a.generation.DataFetched).
					Add(mock.Anything, a.generation.IsCanceled).
					Add(mock.Anything, a.generation.FilesUploaded)
				f.conn.On("Do", args...).Return("", nil)

				args = new(redis.Args).Add("PUBLISH", "generation.updated", a.generation.ID)
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
					Add(mock.Anything, a.generation.DataFetched).
					Add(mock.Anything, a.generation.IsCanceled).
					Add(mock.Anything, a.generation.FilesUploaded).
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
					Add(mock.Anything, a.generation.DataFetched).
					Add(mock.Anything, a.generation.IsCanceled).
					Add(mock.Anything, a.generation.FilesUploaded).
					Add(mock.Anything, a.generation.EndTime.Unix())
				f.conn.On("Do", args...).Return("", nil)

				args = new(redis.Args).Add("PUBLISH", "generation.updated", a.generation.ID)
				f.conn.On("Do", args...).Return("", defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fields := defaultFeedFields()
			tc.setupMocks(tc.args, fields)
			feedRepo := repository.NewFeedRepo(fields.config, fields.client, fields.ftp)

			gotErr := feedRepo.UpdateGenerationState(tc.args.ctx, tc.args.generation)

			assert.Equal(t, tc.wantErr, gotErr)
			fields.assertExpectations(t)
		})
	}
}

func TestFeedRepo_DeleteGeneration(t *testing.T) {
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
				id:  uuid.New().String(),
			},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				f.conn.On("Do", "DEL", a.id).Return("", nil)
			},
		},
		{
			name: "DEL error",
			args: &args{
				ctx: context.Background(),
				id:  uuid.New().String(),
			},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				f.conn.On("Do", "DEL", a.id).Return("", defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fields := defaultFeedFields()
			tc.setupMocks(tc.args, fields)
			feedRepo := repository.NewFeedRepo(fields.config, fields.client, fields.ftp)

			gotErr := feedRepo.DeleteGeneration(tc.args.ctx, tc.args.id)

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
				f.conn.On("Do", "HSET", a.id, "is_canceled", true).Return("", nil)
				f.conn.
					On("Do", "PUBLISH", fmt.Sprintf("%s.canceled", a.id), mock.Anything).
					Return("", nil)
			},
		},
		{
			name: "HSET error",
			args: &args{
				ctx: context.Background(),
				id:  uuid.NewString(),
			},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				f.conn.On("Do", "HSET", a.id, "is_canceled", true).Return("", defaultErr)
			},
			wantErr: defaultErr,
		},
		{
			name: "Publish error",
			args: &args{
				ctx: context.Background(),
				id:  uuid.NewString(),
			},
			setupMocks: func(a *args, f *feedFields) {
				f.client.On("Connection").Return(f.conn)
				f.conn.On("Close").Return(nil)
				f.conn.On("Do", "HSET", a.id, "is_canceled", true).Return("", nil)
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
			feedRepo := repository.NewFeedRepo(fields.config, fields.client, fields.ftp)

			gotErr := feedRepo.CancelGeneration(tc.args.ctx, tc.args.id)

			assert.Equal(t, tc.wantErr, gotErr)
			fields.assertExpectations(t)
		})
	}
}

func TestFeedRepo_OnGenerationCanceled(t *testing.T) {
	type args struct {
		ctx          context.Context
		generationID string
		callback     func()
	}
	testCases := []struct {
		name         string
		args         *args
		setupMocks   func(a *args, f *feedFields)
		mustCallback bool
		wantErr      error
	}{
		{
			name: "succeed",
			args: &args{
				ctx:          context.Background(),
				generationID: uuid.NewString(),
			},
			setupMocks: func(a *args, f *feedFields) {
				channel := fmt.Sprintf("%s.canceled", a.generationID)
				f.client.On("PubSub").Return(f.pubsub)
				f.pubsub.On("Subscribe", channel).Return(nil)

				f.pubsub.On("Receive").Return(redis.Message{Channel: channel, Data: []byte("1")})

				f.pubsub.On("Unsubscribe", channel).Return(nil)
				f.pubsub.On("Close").Return(nil)
			},
			mustCallback: true,
		},
		{
			name: "Subscribe error",
			args: &args{
				ctx:          context.Background(),
				generationID: uuid.NewString(),
			},
			setupMocks: func(a *args, f *feedFields) {
				channel := fmt.Sprintf("%s.canceled", a.generationID)
				f.client.On("PubSub").Return(f.pubsub)
				f.pubsub.On("Subscribe", channel).Return(defaultErr)
				f.pubsub.On("Close").Return(nil)
			},
			mustCallback: false,
			wantErr:      defaultErr,
		},
		{
			name: "Receive error",
			args: &args{
				ctx:          context.Background(),
				generationID: uuid.NewString(),
			},
			setupMocks: func(a *args, f *feedFields) {
				channel := fmt.Sprintf("%s.canceled", a.generationID)
				f.client.On("PubSub").Return(f.pubsub)
				f.pubsub.On("Subscribe", channel).Return(nil)

				f.pubsub.On("Receive").Return(defaultErr)

				f.pubsub.On("Unsubscribe", channel).Return(nil)
				f.pubsub.On("Close").Return(nil)
			},
			mustCallback: false,
			wantErr:      defaultErr,
		},
		{
			name: "context error",
			args: &args{
				ctx:          helper.TimeoutCtx(t, context.Background(), time.Nanosecond),
				generationID: uuid.NewString(),
			},
			setupMocks: func(a *args, f *feedFields) {
				channel := fmt.Sprintf("%s.canceled", a.generationID)
				f.client.On("PubSub").Return(f.pubsub)
				f.pubsub.On("Subscribe", channel).Return(nil)

				f.pubsub.On("Receive").
					Return(redis.Message{Channel: channel, Data: []byte("1")}).
					After(time.Second).Maybe()

				f.pubsub.On("Unsubscribe", channel).Return(nil)
				f.pubsub.On("Close").Return(nil)
			},
			mustCallback: false,
			wantErr:      context.DeadlineExceeded,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var callbackCalled bool
			fields := defaultFeedFields()
			tc.setupMocks(tc.args, fields)
			feedRepo := repository.NewFeedRepo(fields.config, fields.client, fields.ftp)
			tc.args.callback = func() {
				callbackCalled = true
			}

			gotErr := feedRepo.OnGenerationCanceled(tc.args.ctx, tc.args.generationID, tc.args.callback)

			assert.Equal(t, tc.wantErr, gotErr)
			assert.Equal(t, tc.mustCallback, callbackCalled, "Callback called mismatch")
			fields.assertExpectations(t)
		})
	}
}

func TestFeedRepo_OnGenerationsUpdated(t *testing.T) {
	type args struct {
		ctx      context.Context
		callback func(generation *entity.Generation)
	}
	testCases := []struct {
		name       string
		args       *args
		setupMocks func(a *args, f *feedFields)
		want       *entity.Generation
		wantErr    error
	}{
		{
			name: "succeed",
			args: &args{ctx: helper.TimeoutCtx(t, context.Background(), time.Millisecond)},
			setupMocks: func(a *args, f *feedFields) {
				channel := "generation.updated"
				f.client.On("PubSub").Return(f.pubsub)
				f.client.On("Connection").Return(f.conn)
				f.pubsub.On("Subscribe", channel).Return(nil)

				f.pubsub.On("Receive").Return(redis.Message{Channel: channel, Data: []byte("abc")})
				f.conn.On("Do", "HGETALL", "abc").Return([]interface{}{
					[]byte("type"), []byte("test"),
					[]byte("progress"), []byte("99"),
					[]byte("files_uploaded"), []byte("5"),
					[]byte("data_fetched"), []byte("1"),
					[]byte("start_time"), []byte(strconv.Itoa(int(time.Unix(10, 0).Unix()))),
				}, nil)

				f.pubsub.On("Unsubscribe", channel).Return(nil)
				f.pubsub.On("Close").Return(nil)
				f.conn.On("Close").Return(nil)
			},
			want: &entity.Generation{
				ID:            "abc",
				Type:          "test",
				Progress:      99,
				DataFetched:   true,
				FilesUploaded: 5,
				StartTime:     time.Unix(10, 0),
			},
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "Subscribe error",
			args: &args{ctx: context.Background()},
			setupMocks: func(a *args, f *feedFields) {
				channel := "generation.updated"
				f.client.On("PubSub").Return(f.pubsub)
				f.pubsub.On("Subscribe", channel).Return(defaultErr)
				f.pubsub.On("Close").Return(nil)
			},
			wantErr: defaultErr,
		},
		{
			name: "Receive error",
			args: &args{
				ctx: context.Background(),
			},
			setupMocks: func(a *args, f *feedFields) {
				channel := "generation.updated"
				f.client.On("PubSub").Return(f.pubsub)
				f.pubsub.On("Subscribe", channel).Return(nil)

				f.pubsub.On("Receive").Return(defaultErr)

				f.pubsub.On("Unsubscribe", channel).Return(nil)
				f.pubsub.On("Close").Return(nil)
			},
			wantErr: defaultErr,
		},
		{
			name: "context error",
			args: &args{
				ctx: helper.TimeoutCtx(t, context.Background(), time.Nanosecond),
			},
			setupMocks: func(a *args, f *feedFields) {
				channel := "generation.updated"
				f.client.On("Connection").Return(f.conn).Maybe()
				f.client.On("PubSub").Return(f.pubsub)
				f.pubsub.On("Subscribe", channel).Return(nil)

				f.pubsub.On("Receive").
					Return(redis.Message{Channel: channel, Data: []byte("1")}).
					After(time.Second).Maybe()
				f.conn.On("Do", "HGETALL", "1").Return([]interface{}{
					[]byte("type"), []byte("test"),
				}, nil).Maybe()

				f.pubsub.On("Unsubscribe", channel).Return(nil)
				f.pubsub.On("Close").Return(nil)
				f.conn.On("Close").Return(nil).Maybe()
			},
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "HGETALL error",
			args: &args{ctx: helper.TimeoutCtx(t, context.Background(), time.Second)},
			setupMocks: func(a *args, f *feedFields) {
				channel := "generation.updated"
				f.client.On("PubSub").Return(f.pubsub)
				f.client.On("Connection").Return(f.conn)
				f.pubsub.On("Subscribe", channel).Return(nil)

				f.pubsub.On("Receive").Return(redis.Message{Channel: channel, Data: []byte("abc")})
				f.conn.On("Do", "HGETALL", "abc").Return(nil, defaultErr)

				f.pubsub.On("Unsubscribe", channel).Return(nil)
				f.pubsub.On("Close").Return(nil)
				f.conn.On("Close").Return(nil)
			},
			wantErr: defaultErr,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got *entity.Generation
			fields := defaultFeedFields()
			tc.setupMocks(tc.args, fields)
			feedRepo := repository.NewFeedRepo(fields.config, fields.client, fields.ftp)
			tc.args.callback = func(g *entity.Generation) {
				got = g
			}

			gotErr := feedRepo.OnGenerationsUpdated(tc.args.ctx, tc.args.callback)

			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErr, gotErr)
			fields.assertExpectations(t)
		})
	}
}
