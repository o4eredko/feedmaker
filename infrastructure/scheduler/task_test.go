package scheduler_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-feedmaker/infrastructure/scheduler"
	"go-feedmaker/infrastructure/scheduler/mocks"
)

type (
	taskFields struct {
		cmd      *mocks.Runner
		schedule *scheduler.Schedule
	}
)

func defaultTaskFields() *taskFields {
	return &taskFields{
		cmd:      new(mocks.Runner),
		schedule: scheduler.NewSchedule(time.Now().UTC(), time.Second*42),
	}
}

func TestNewTask(t *testing.T) {
	fields := defaultTaskFields()
	gotTask := scheduler.NewTask(fields.cmd, fields.schedule)
	assert.Equal(t, fields.cmd, gotTask.Cmd)
	assert.Equal(t, fields.schedule, gotTask.Schedule)
}
