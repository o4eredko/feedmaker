package scheduler_test

import (
	"strconv"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"

	"go-feedmaker/adapter/repository/mocks"
	"go-feedmaker/infrastructure/scheduler"
)

type (
	scheduleSaverFields struct {
		client *mocks.RedisClient
		conn   *mocks.Connection
	}
)

func defaultScheduleSaverFields() *scheduleSaverFields {
	return &scheduleSaverFields{
		client: new(mocks.RedisClient),
		conn:   new(mocks.Connection),
	}
}

func TestNewScheduleSaver(t *testing.T) {
	fields := defaultScheduleSaverFields()
	saver := scheduler.NewScheduleSaver(fields.client)
	assert.Equal(t, fields.client, saver.RedisClient())
}

func Test_scheduleSaver_Store(t *testing.T) {
	type args struct {
		id       scheduler.TaskID
		schedule *scheduler.Schedule
	}
	defaultArgs := func() *args {
		return &args{
			id:       defaultTaskID,
			schedule: defaultSchedule,
		}
	}
	testCases := []struct {
		name       string
		fields     *scheduleSaverFields
		args       *args
		setupMocks func(*scheduleSaverFields, *args)
		wantErr    error
	}{
		{
			name:   "succeed",
			fields: defaultScheduleSaverFields(),
			args:   defaultArgs(),
			setupMocks: func(fields *scheduleSaverFields, args *args) {
				fields.client.On("Connection").Return(fields.conn)
				fields.conn.On("Close").Return(nil)
				fields.conn.On("Send", "MULTI").Return(nil)
				fields.conn.On("Send", "SADD", scheduler.TaskIDsKey, args.id).Return(nil)
				argsToSend := new(redis.Args).
					Add("HMSET", args.id).
					Add("start_timestamp", args.schedule.StartTimestamp().Unix()).
					Add("fire_interval", args.schedule.FireInterval().Seconds())
				fields.conn.On("Send", argsToSend...).Return(nil)
				fields.conn.On("Do", "EXEC").Return("OK", nil)
			},
		},
		{
			name:   "conn.Do returns error",
			fields: defaultScheduleSaverFields(),
			args:   defaultArgs(),
			setupMocks: func(fields *scheduleSaverFields, args *args) {
				fields.client.On("Connection").Return(fields.conn)
				fields.conn.On("Close").Return(nil)
				fields.conn.On("Send", "MULTI").Return(nil)
				fields.conn.On("Send", "SADD", scheduler.TaskIDsKey, args.id).Return(nil)
				argsToSend := new(redis.Args).
					Add("HMSET", args.id).
					Add("start_timestamp", args.schedule.StartTimestamp().Unix()).
					Add("fire_interval", args.schedule.FireInterval().Seconds())
				fields.conn.On("Send", argsToSend...).Return(nil)
				fields.conn.On("Do", "EXEC").Return("", defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields, testCase.args)
			s := scheduler.NewScheduleSaver(testCase.fields.client)
			gotErr := s.Store(testCase.args.id, testCase.args.schedule)
			assert.Equal(t, testCase.wantErr, gotErr)
			testCase.fields.conn.AssertExpectations(t)
			testCase.fields.client.AssertExpectations(t)
		})
	}
}

func Test_scheduleSaver_Load(t *testing.T) {
	type args struct {
		id scheduler.TaskID
	}
	defaultArgs := func() *args {
		return &args{
			id: defaultTaskID,
		}
	}
	testCases := []struct {
		name         string
		fields       *scheduleSaverFields
		args         *args
		setupMocks   func(*scheduleSaverFields, *args)
		wantSchedule *scheduler.Schedule
		wantErr      error
	}{
		{
			name:   "succeed",
			fields: defaultScheduleSaverFields(),
			args:   defaultArgs(),
			setupMocks: func(fields *scheduleSaverFields, args *args) {
				fields.client.On("Connection").Return(fields.conn)
				fields.conn.On("Close").Return(nil)
				rawSchedule := []interface{}{
					[]byte("start_timestamp"), []byte(strconv.Itoa(int(defaultSchedule.StartTimestamp().Unix()))),
					[]byte("fire_interval"), []byte(strconv.Itoa(int(defaultSchedule.FireInterval().Seconds()))),
				}
				fields.conn.On("Do", "HGETALL", args.id).Return(rawSchedule, nil)
			},
			wantSchedule: defaultSchedule,
		},
		{
			name:   "conn.Do returns error",
			fields: defaultScheduleSaverFields(),
			args:   defaultArgs(),
			setupMocks: func(fields *scheduleSaverFields, args *args) {
				fields.client.On("Connection").Return(fields.conn)
				fields.conn.On("Close").Return(nil)
				fields.conn.On("Do", "HGETALL", args.id).Return(nil, defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields, testCase.args)
			s := scheduler.NewScheduleSaver(testCase.fields.client)
			gotSchedule, gotErr := s.Load(testCase.args.id)
			assert.Equal(t, testCase.wantErr, gotErr)
			assert.Equal(t, testCase.wantSchedule, gotSchedule)
			testCase.fields.conn.AssertExpectations(t)
			testCase.fields.client.AssertExpectations(t)
		})
	}
}

func Test_scheduleSaver_Delete(t *testing.T) {
	type args struct {
		id scheduler.TaskID
	}
	defaultArgs := func() *args {
		return &args{
			id: defaultTaskID,
		}
	}
	testCases := []struct {
		name       string
		fields     *scheduleSaverFields
		args       *args
		setupMocks func(*scheduleSaverFields, *args)
		wantErr    error
	}{
		{
			name:   "succeed",
			fields: defaultScheduleSaverFields(),
			args:   defaultArgs(),
			setupMocks: func(fields *scheduleSaverFields, args *args) {
				fields.client.On("Connection").Return(fields.conn)
				fields.conn.On("Close").Return(nil)
				fields.conn.On("Send", "MULTI").Return(nil)
				fields.conn.On("Send", "SREM", scheduler.TaskIDsKey, args.id).Return(nil)
				fields.conn.On("Send", "DEL", args.id).Return(nil)
				fields.conn.On("Do", "EXEC").Return("OK", nil)
			},
		},
		{
			name:   "conn.Do returns error",
			fields: defaultScheduleSaverFields(),
			args:   defaultArgs(),
			setupMocks: func(fields *scheduleSaverFields, args *args) {
				fields.client.On("Connection").Return(fields.conn)
				fields.conn.On("Close").Return(nil)
				fields.conn.On("Send", "MULTI").Return(nil)
				fields.conn.On("Send", "SREM", scheduler.TaskIDsKey, args.id).Return(nil)
				fields.conn.On("Send", "DEL", args.id).Return(nil)
				fields.conn.On("Do", "EXEC").Return("", defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields, testCase.args)
			s := scheduler.NewScheduleSaver(testCase.fields.client)
			gotErr := s.Delete(testCase.args.id)
			assert.Equal(t, testCase.wantErr, gotErr)
			testCase.fields.conn.AssertExpectations(t)
			testCase.fields.client.AssertExpectations(t)
		})
	}
}

func Test_scheduleSaver_ListScheduledTaskIDs(t *testing.T) {
	testCases := []struct {
		name        string
		fields      *scheduleSaverFields
		setupMocks  func(*scheduleSaverFields)
		wantTaskIDs []scheduler.TaskID
		wantErr     error
	}{
		{
			name:   "succeed",
			fields: defaultScheduleSaverFields(),
			setupMocks: func(fields *scheduleSaverFields) {
				fields.client.On("Connection").Return(fields.conn)
				fields.conn.On("Close").Return(nil)
				ids := make([]interface{}, len(defaultScheduledTaskIDs))
				for i, taskID := range defaultScheduledTaskIDs {
					ids[i] = string(taskID)
				}
				fields.conn.
					On("Do", "SMEMBERS", scheduler.TaskIDsKey).
					Return(ids, nil)
			},
			wantTaskIDs: defaultScheduledTaskIDs,
		},
		{
			name:   "conn.Do returns error",
			fields: defaultScheduleSaverFields(),
			setupMocks: func(fields *scheduleSaverFields) {
				fields.client.On("Connection").Return(fields.conn)
				fields.conn.On("Close").Return(nil)
				fields.conn.On("Do", "SMEMBERS", scheduler.TaskIDsKey).Return(nil, defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields)
			s := scheduler.NewScheduleSaver(testCase.fields.client)
			gotTaskIDs, gotErr := s.ListScheduledTaskIDs()
			assert.Equal(t, testCase.wantErr, gotErr)
			assert.Equal(t, testCase.wantTaskIDs, gotTaskIDs)
			testCase.fields.conn.AssertExpectations(t)
			testCase.fields.client.AssertExpectations(t)
		})
	}
}
