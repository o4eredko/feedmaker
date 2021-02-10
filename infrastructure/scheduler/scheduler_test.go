package scheduler_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-feedmaker/infrastructure/scheduler"
	"go-feedmaker/infrastructure/scheduler/mocks"
)

type (
	schedulerFields struct {
		cron   *mocks.Croner
		mapper *mocks.TaskIDMapper
	}
)

func defaultSchedulerFields() *schedulerFields {
	return &schedulerFields{
		cron:   new(mocks.Croner),
		mapper: new(mocks.TaskIDMapper),
	}
}

func TestNew(t *testing.T) {
	fields := defaultSchedulerFields()
	s := scheduler.New(fields.cron, fields.mapper)
	assert.Equal(t, fields.cron, s.Cron())
	assert.Equal(t, fields.mapper, s.Mapper())
}

func TestScheduler_Start(t *testing.T) {
	testsCases := []struct {
		name       string
		fields     *schedulerFields
		setupMocks func(*schedulerFields)
	}{
		{
			name:   "succeed",
			fields: defaultSchedulerFields(),
			setupMocks: func(fields *schedulerFields) {
				fields.cron.On("Start")
			},
		},
	}
	for _, testsCase := range testsCases {
		t.Run(testsCase.name, func(t *testing.T) {
			testsCase.setupMocks(testsCase.fields)
			s := scheduler.New(testsCase.fields.cron, testsCase.fields.mapper)
			s.Start()

			testsCase.fields.cron.AssertExpectations(t)
		})
	}
}

func TestScheduler_Stop(t *testing.T) {
	testsCases := []struct {
		name       string
		fields     *schedulerFields
		setupMocks func(*schedulerFields)
	}{
		{
			name:   "succeed",
			fields: defaultSchedulerFields(),
			setupMocks: func(fields *schedulerFields) {
				fields.cron.
					On("Stop").
					Return(context.TODO())
			},
		},
	}
	for _, testsCase := range testsCases {
		t.Run(testsCase.name, func(t *testing.T) {
			testsCase.setupMocks(testsCase.fields)
			s := scheduler.New(testsCase.fields.cron, testsCase.fields.mapper)
			s.Stop()

			testsCase.fields.cron.AssertExpectations(t)
		})
	}
}

func TestScheduler_AddTask(t *testing.T) {
	type args struct {
		taskID scheduler.TaskID
		task   scheduler.Task
	}
	testsCases := []struct {
		name       string
		fields     *schedulerFields
		setupMocks func(*schedulerFields, *args)
		args       *args
		wantErr    error
	}{
		{
			name:   "succeed",
			fields: defaultSchedulerFields(),
			setupMocks: func(fields *schedulerFields, args *args) {
				fields.mapper.
					On("Load", args.taskID).
					Return(defaultEntryID, nil)
				fields.cron.
					On("Schedule", args.task.Schedule(), args.task.Cmd()).
					Return(defaultEntryID)
				fields.mapper.
					On("Store", args.taskID, defaultEntryID).
					Return(nil)
			},
			args: &args{
				taskID: defaultTaskID,
				task:   defaultTask,
			},
		},
		{
			name:   "loading error",
			fields: defaultSchedulerFields(),
			setupMocks: func(fields *schedulerFields, args *args) {
				fields.mapper.
					On("Load", args.taskID).
					Return(defaultEntryID, defaultErr)
			},
			args: &args{
				taskID: defaultTaskID,
				task:   defaultTask,
			},
			wantErr: defaultErr,
		},
		{
			name:   "storing error",
			fields: defaultSchedulerFields(),
			setupMocks: func(fields *schedulerFields, args *args) {
				fields.mapper.
					On("Load", args.taskID).
					Return(defaultEntryID, nil)
				fields.cron.
					On("Schedule", args.task.Schedule(), args.task.Cmd()).
					Return(defaultEntryID)
				fields.mapper.
					On("Store", args.taskID, defaultEntryID).
					Return(defaultErr)
			},
			args: &args{
				taskID: defaultTaskID,
				task:   defaultTask,
			},
			wantErr: defaultErr,
		},
	}
	for _, testsCase := range testsCases {
		t.Run(testsCase.name, func(t *testing.T) {
			testsCase.setupMocks(testsCase.fields, testsCase.args)
			s := scheduler.New(testsCase.fields.cron, testsCase.fields.mapper)
			gotErr := s.AddTask(testsCase.args.taskID, testsCase.args.task)
			assert.Equal(t, testsCase.wantErr, gotErr)

			testsCase.fields.cron.AssertExpectations(t)
			testsCase.fields.mapper.AssertExpectations(t)
		})
	}
}

func TestScheduler_RemoveTask(t *testing.T) {
	type args struct {
		taskID scheduler.TaskID
	}
	testsCases := []struct {
		name       string
		fields     *schedulerFields
		setupMocks func(*schedulerFields, *args)
		args       *args
		wantErr    error
	}{
		{
			name:   "succeed",
			fields: defaultSchedulerFields(),
			setupMocks: func(fields *schedulerFields, args *args) {
				fields.mapper.
					On("Load", args.taskID).
					Return(defaultEntryID, nil)
				fields.cron.
					On("Remove", defaultEntryID)
				fields.mapper.
					On("Delete", args.taskID).
					Return(nil)
			},
			args: &args{
				taskID: defaultTaskID,
			},
		},
		{
			name:   "loading error",
			fields: defaultSchedulerFields(),
			setupMocks: func(fields *schedulerFields, args *args) {
				fields.mapper.
					On("Load", args.taskID).
					Return(defaultEntryID, defaultErr)
			},
			args: &args{
				taskID: defaultTaskID,
			},
			wantErr: defaultErr,
		},
		{
			name:   "deleting error",
			fields: defaultSchedulerFields(),
			setupMocks: func(fields *schedulerFields, args *args) {
				fields.mapper.
					On("Load", args.taskID).
					Return(defaultEntryID, nil)
				fields.cron.
					On("Remove", defaultEntryID)
				fields.mapper.
					On("Delete", args.taskID).
					Return(defaultErr)
			},
			args: &args{
				taskID: defaultTaskID,
			},
			wantErr: defaultErr,
		},
	}
	for _, testsCase := range testsCases {
		t.Run(testsCase.name, func(t *testing.T) {
			testsCase.setupMocks(testsCase.fields, testsCase.args)
			s := scheduler.New(testsCase.fields.cron, testsCase.fields.mapper)
			gotErr := s.RemoveTask(testsCase.args.taskID)
			assert.Equal(t, testsCase.wantErr, gotErr)

			testsCase.fields.cron.AssertExpectations(t)
			testsCase.fields.mapper.AssertExpectations(t)
		})
	}
}
