package scheduler_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-feedmaker/infrastructure/scheduler"
	interactorMocks "go-feedmaker/interactor/mocks"
)

func TestScheduler_ScheduleGeneration(t *testing.T) {
	type args struct {
		feedsInteractor *interactorMocks.FeedInteractor
		taskID          scheduler.TaskID
		schedule        *scheduler.Schedule
	}
	defaultArgs := func() *args {
		return &args{
			feedsInteractor: new(interactorMocks.FeedInteractor),
			taskID:          defaultTaskID,
			schedule:        defaultSchedule,
		}
	}
	makeCmdMatches := func(args *args) func(cmd *scheduler.Cmd) bool {
		return func(cmd *scheduler.Cmd) bool {
			return cmd.Func.Type() == reflect.TypeOf(args.feedsInteractor.GenerateFeed) &&
				cmd.Args[0].Type() == reflect.TypeOf(context.Background()) &&
				cmd.Args[1].Kind() == reflect.String &&
				cmd.Args[1].String() == reflect.ValueOf(args.taskID).String()
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
					On("Store", args.taskID, args.schedule).
					Return(nil)
				cmdMatches := makeCmdMatches(args)
				fields.cron.
					On("Schedule", args.schedule, mock.MatchedBy(cmdMatches)).
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
					On("Store", args.taskID, args.schedule).
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
					On("Store", args.taskID, args.schedule).
					Return(nil)
				cmdMatches := makeCmdMatches(args)
				fields.cron.
					On("Schedule", args.schedule, mock.MatchedBy(cmdMatches)).
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
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields, testCase.args)
			s := scheduler.New(testCase.fields.cron, testCase.fields.saver)
			s.SetMapper(testCase.fields.mapper)

			gotErr := s.ScheduleGeneration(testCase.args.feedsInteractor,
				string(testCase.args.taskID), testCase.args.schedule)
			assert.Equal(t, testCase.wantErr, gotErr)

			testCase.fields.cron.AssertExpectations(t)
			testCase.fields.saver.AssertExpectations(t)
			testCase.fields.mapper.AssertExpectations(t)
		})
	}
}
