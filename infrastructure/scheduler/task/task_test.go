package task_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-feedmaker/infrastructure/scheduler/mocks"
	"go-feedmaker/infrastructure/scheduler/task"
)

type (
	taskFields struct {
		cmd      *mocks.Runner
		schedule *task.Schedule
	}
)

func defaultTaskFields() *taskFields {
	return &taskFields{
		cmd:      new(mocks.Runner),
		schedule: task.NewSchedule(time.Now().UTC(), time.Second*42),
	}
}

func TestNew(t *testing.T) {
	fields := defaultTaskFields()
	gotTask := task.NewTask(fields.cmd, fields.schedule)
	assert.Equal(t, fields.cmd, gotTask.Cmd)
	assert.Equal(t, fields.schedule, gotTask.Schedule)
}
