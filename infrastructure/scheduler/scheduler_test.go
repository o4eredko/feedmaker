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
		saver  *mocks.ScheduleSaver
	}
)

func defaultSchedulerFields() *schedulerFields {
	return &schedulerFields{
		cron:   new(mocks.Croner),
		mapper: new(mocks.TaskIDMapper),
		saver:  new(mocks.ScheduleSaver),
	}
}

func TestNew(t *testing.T) {
	fields := defaultSchedulerFields()
	s := scheduler.New(fields.cron, fields.saver)
	s.SetMapper(fields.mapper)
	assert.Equal(t, fields.cron, s.Cron())
	assert.Equal(t, fields.saver, s.Saver())
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
			s := scheduler.New(testsCase.fields.cron, testsCase.fields.saver)
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
				fields.cron.On("Stop").Return(context.TODO())
			},
		},
	}
	for _, testsCase := range testsCases {
		t.Run(testsCase.name, func(t *testing.T) {
			testsCase.setupMocks(testsCase.fields)
			s := scheduler.New(testsCase.fields.cron, testsCase.fields.saver)
			s.Stop()

			testsCase.fields.cron.AssertExpectations(t)
		})
	}
}

func TestScheduler_ScheduleTask(t *testing.T) {
	type args struct {
		taskID scheduler.TaskID
		task   *scheduler.Task
	}
	defaultArgs := func() *args {
		cmd := new(mocks.Runner)
		schedule := defaultSchedule
		return &args{
			taskID: defaultTaskID,
			task:   scheduler.NewTask(cmd, schedule),
		}
	}
	testCases := []struct {
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
				fields.saver.
					On("Store", args.taskID, args.task.Schedule).
					Return(nil)
				fields.cron.
					On("Schedule", args.task.Schedule, args.task.Cmd).
					Return(defaultEntryID)
				fields.mapper.
					On("Store", args.taskID, defaultEntryID).
					Return(nil)
			},
			args: defaultArgs(),
		},
		{
			name:   "saver.Store returns error",
			fields: defaultSchedulerFields(),
			setupMocks: func(fields *schedulerFields, args *args) {
				fields.saver.
					On("Store", args.taskID, args.task.Schedule).
					Return(defaultErr)
			},
			args:    defaultArgs(),
			wantErr: defaultErr,
		},
		{
			name:   "mapper.Store returns error",
			fields: defaultSchedulerFields(),
			setupMocks: func(fields *schedulerFields, args *args) {
				fields.saver.
					On("Store", args.taskID, args.task.Schedule).
					Return(nil)
				fields.cron.
					On("Schedule", args.task.Schedule, args.task.Cmd).
					Return(defaultEntryID)
				fields.mapper.
					On("Store", args.taskID, defaultEntryID).
					Return(defaultErr)
				fields.cron.
					On("Remove", defaultEntryID)
			},
			args:    defaultArgs(),
			wantErr: defaultErr,
		},
	}
	for _, testsCase := range testCases {
		t.Run(testsCase.name, func(t *testing.T) {
			testsCase.setupMocks(testsCase.fields, testsCase.args)
			s := scheduler.New(testsCase.fields.cron, testsCase.fields.saver)
			s.SetMapper(testsCase.fields.mapper)

			gotErr := s.ScheduleTask(testsCase.args.taskID, testsCase.args.task)
			assert.Equal(t, testsCase.wantErr, gotErr)

			testsCase.fields.cron.AssertExpectations(t)
			testsCase.fields.saver.AssertExpectations(t)
			testsCase.fields.mapper.AssertExpectations(t)
		})
	}
}

func TestScheduler_RemoveTask(t *testing.T) {
	type args struct {
		taskID scheduler.TaskID
	}
	defaultArgs := func() *args {
		return &args{
			taskID: defaultTaskID,
		}
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
				fields.saver.
					On("Delete", args.taskID).
					Return(nil)
			},
			args: defaultArgs(),
		},
		{
			name:   "mapper.Load returns error",
			fields: defaultSchedulerFields(),
			setupMocks: func(fields *schedulerFields, args *args) {
				fields.mapper.
					On("Load", args.taskID).
					Return(defaultEntryID, defaultErr)
			},
			args:    defaultArgs(),
			wantErr: defaultErr,
		},
		{
			name:   "mapper.Delete returns error",
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
			args:    defaultArgs(),
			wantErr: defaultErr,
		},
		{
			name:   "saver.Delete returns error",
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
				fields.saver.
					On("Delete", args.taskID).
					Return(defaultErr)
			},
			args:    defaultArgs(),
			wantErr: defaultErr,
		},
	}
	for _, testsCase := range testsCases {
		t.Run(testsCase.name, func(t *testing.T) {
			testsCase.setupMocks(testsCase.fields, testsCase.args)
			s := scheduler.New(testsCase.fields.cron, testsCase.fields.saver)
			s.SetMapper(testsCase.fields.mapper)

			gotErr := s.RemoveTask(testsCase.args.taskID)
			assert.Equal(t, testsCase.wantErr, gotErr)

			testsCase.fields.cron.AssertExpectations(t)
			testsCase.fields.saver.AssertExpectations(t)
			testsCase.fields.mapper.AssertExpectations(t)
		})
	}
}

func TestScheduler_ListSchedules(t *testing.T) {
	testCases := []struct {
		name          string
		fields        *schedulerFields
		setupMocks    func(*schedulerFields)
		wantSchedules map[scheduler.TaskID]*scheduler.Schedule
		wantErr       error
	}{
		{
			name:   "succeed",
			fields: defaultSchedulerFields(),
			setupMocks: func(fields *schedulerFields) {
				fields.saver.
					On("ListScheduledTaskIDs").
					Return(defaultScheduledTaskIDs, nil)
				for id, schedule := range defaultTaskSchedules {
					fields.saver.
						On("Load", id).
						Return(schedule, nil)
				}
			},
			wantSchedules: defaultTaskSchedules,
		},
		{
			name:   "saver.ListScheduledTaskIDs returns error",
			fields: defaultSchedulerFields(),
			setupMocks: func(fields *schedulerFields) {
				fields.saver.
					On("ListScheduledTaskIDs").
					Return(nil, defaultErr)
			},
			wantErr: defaultErr,
		},
		{
			name:   "saver.Load returns error",
			fields: defaultSchedulerFields(),
			setupMocks: func(fields *schedulerFields) {
				fields.saver.
					On("ListScheduledTaskIDs").
					Return(defaultScheduledTaskIDs, nil)
				fields.saver.
					On("Load", defaultScheduledTaskIDs[0]).
					Return(nil, defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields)
			s := scheduler.New(testCase.fields.cron, testCase.fields.saver)
			s.SetMapper(testCase.fields.mapper)

			gotSchedules, gotErr := s.ListSchedules()
			assert.Equal(t, testCase.wantSchedules, gotSchedules)
			assert.Equal(t, testCase.wantErr, gotErr)

			testCase.fields.cron.AssertExpectations(t)
			testCase.fields.saver.AssertExpectations(t)
			testCase.fields.mapper.AssertExpectations(t)
		})
	}
}
