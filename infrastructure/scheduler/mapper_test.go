package scheduler_test

import (
	"testing"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"

	"go-feedmaker/infrastructure/scheduler"
)

type (
	mapperFields struct {
		mapping map[scheduler.TaskID]cron.EntryID
	}
)

func TestMapper_Store(t *testing.T) {
	type args struct {
		taskID  scheduler.TaskID
		entryID cron.EntryID
	}
	testCases := []struct {
		name    string
		fields  *mapperFields
		args    args
		wantErr error
	}{
		{
			name: "succeed",
			fields: &mapperFields{
				mapping: map[scheduler.TaskID]cron.EntryID{},
			},
			args: args{
				taskID:  defaultTaskID,
				entryID: defaultEntryID,
			},
		},
		{
			name: "task id duplication error",
			fields: &mapperFields{
				mapping: map[scheduler.TaskID]cron.EntryID{
					defaultTaskID: defaultEntryID,
				},
			},
			args: args{
				taskID:  defaultTaskID,
				entryID: defaultEntryID,
			},
			wantErr: scheduler.ErrTaskAlreadyExists,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			m := scheduler.NewMapper()
			m.SetMapping(testCase.fields.mapping)

			gotErr := m.Store(testCase.args.taskID, testCase.args.entryID)
			assert.Equal(t, testCase.wantErr, gotErr)
		})
	}
}

func TestMapper_Load(t *testing.T) {
	type args struct {
		taskID scheduler.TaskID
	}
	testCases := []struct {
		name        string
		fields      *mapperFields
		args        args
		wantEntryID cron.EntryID
		wantErr     error
	}{
		{
			name: "succeed",
			fields: &mapperFields{
				mapping: map[scheduler.TaskID]cron.EntryID{
					defaultTaskID: defaultEntryID,
				},
			},
			args: args{
				taskID: defaultTaskID,
			},
			wantEntryID: defaultEntryID,
		},
		{
			name: "task id lookup error",
			fields: &mapperFields{
				mapping: map[scheduler.TaskID]cron.EntryID{},
			},
			args: args{
				taskID: defaultTaskID,
			},
			wantErr: scheduler.ErrTaskNotFound,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			m := scheduler.NewMapper()
			m.SetMapping(testCase.fields.mapping)

			gotEntryID, gotErr := m.Load(testCase.args.taskID)
			assert.Equal(t, testCase.wantEntryID, gotEntryID)
			assert.Equal(t, testCase.wantErr, gotErr)
		})
	}
}

func TestMapper_Delete(t *testing.T) {
	type args struct {
		taskID scheduler.TaskID
	}
	testCases := []struct {
		name    string
		fields  *mapperFields
		args    args
		wantErr error
	}{
		{
			name: "succeed",
			fields: &mapperFields{
				mapping: map[scheduler.TaskID]cron.EntryID{
					defaultTaskID: defaultEntryID,
				},
			},
			args: args{
				taskID: defaultTaskID,
			},
		},
		{
			name: "task id lookup error",
			fields: &mapperFields{
				mapping: map[scheduler.TaskID]cron.EntryID{},
			},
			args: args{
				taskID: defaultTaskID,
			},
			wantErr: scheduler.ErrTaskNotFound,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			m := scheduler.NewMapper()
			m.SetMapping(testCase.fields.mapping)

			gotErr := m.Delete(testCase.args.taskID)
			assert.Equal(t, testCase.wantErr, gotErr)
		})
	}
}
