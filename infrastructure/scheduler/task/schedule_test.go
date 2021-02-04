package task_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-feedmaker/infrastructure/scheduler/task"
)

type (
	scheduleWithSpecificStartTimestampFields struct {
		startTimestamp time.Time
		delayInterval  time.Duration
	}
)

var (
	defaultStartTimestamp = time.Now().UTC().Add(time.Hour * 13)
	defaultDelayInterval  = time.Hour * 42
)

func defaultScheduleWithSpecificStartTimestampFields() *scheduleWithSpecificStartTimestampFields {
	return &scheduleWithSpecificStartTimestampFields{
		startTimestamp: defaultStartTimestamp,
		delayInterval:  defaultDelayInterval,
	}
}

func getNowTimestamp() time.Time {
	now := time.Now()
	return now.UTC()
}

func TestNewScheduleWithSpecificStartTimestamp(t *testing.T) {
	fields := defaultScheduleWithSpecificStartTimestampFields()
	s := task.NewScheduleWithSpecificStartTimestamp(fields.startTimestamp, fields.delayInterval)
	assert.Equal(t, fields.startTimestamp, s.StartTimestamp)
	assert.Equal(t, fields.delayInterval, s.DelayInterval)
}

func TestScheduleWithSpecificStartTimestamp_Next(t *testing.T) {
	testCases := []struct {
		name   string
		fields *scheduleWithSpecificStartTimestampFields
	}{
		{
			name:   "succeed",
			fields: defaultScheduleWithSpecificStartTimestampFields(),
		},
		{
			name: "succeed with start timestamp alignment",
			fields: &scheduleWithSpecificStartTimestampFields{
				startTimestamp: time.Unix(0, 0),
				delayInterval:  defaultDelayInterval,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			s := &task.ScheduleWithSpecificStartTimestamp{
				StartTimestamp: testCase.fields.startTimestamp,
				DelayInterval:  testCase.fields.delayInterval,
			}

			nowTimestamp := getNowTimestamp()
			elapsed := nowTimestamp.Sub(testCase.fields.startTimestamp)
			expectedStartAt := testCase.fields.startTimestamp.Add(elapsed.Round(testCase.fields.delayInterval))
			gotStartAt := s.Next(nowTimestamp)
			assert.Equal(t, expectedStartAt, gotStartAt)

			nowTimestamp = getNowTimestamp()
			expectedNext := nowTimestamp.Add(testCase.fields.delayInterval).Truncate(time.Second)
			gotNext := s.Next(nowTimestamp)
			assert.Equal(t, expectedNext, gotNext)
		})
	}
}
