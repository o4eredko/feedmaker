package scheduler_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-feedmaker/infrastructure/scheduler"
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
	s := scheduler.NewSchedule(fields.startTimestamp, fields.delayInterval)
	assert.Equal(t, fields.startTimestamp, s.StartTimestamp())
	assert.Equal(t, fields.delayInterval, s.FireInterval())
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
			s := scheduler.NewSchedule(testCase.fields.startTimestamp, testCase.fields.delayInterval)

			nowTimestamp := getNowTimestamp()
			elapsedTimeRoundedToInterval := nowTimestamp.Sub(testCase.fields.startTimestamp).Round(testCase.fields.delayInterval)
			expectedStartAt := testCase.fields.startTimestamp.Add(elapsedTimeRoundedToInterval)
			if expectedStartAt.Before(nowTimestamp) {
				expectedStartAt = expectedStartAt.Add(testCase.fields.delayInterval)
			}
			gotStartAt := s.Next(nowTimestamp)
			assert.Equal(t, expectedStartAt, gotStartAt)

			nowTimestamp = getNowTimestamp()
			expectedNext := nowTimestamp.Add(testCase.fields.delayInterval).Truncate(time.Second)
			gotNext := s.Next(nowTimestamp)
			assert.Equal(t, expectedNext, gotNext)
		})
	}
}
