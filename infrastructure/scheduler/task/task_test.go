package task_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-feedmaker/infrastructure/scheduler/mocks"
	"go-feedmaker/infrastructure/scheduler/task"
)

type (
	taskFields struct {
		cmd      *mocks.Runner
		schedule *mocks.Nexter
	}
)

func defaultTaskFields() *taskFields {
	return &taskFields{
		cmd:      new(mocks.Runner),
		schedule: new(mocks.Nexter),
	}
}

func TestNew(t *testing.T) {
	fields := defaultTaskFields()
	gotTask := task.New(fields.cmd, fields.schedule)
	assert.Equal(t, fields.cmd, gotTask.Cmd())
	assert.Equal(t, fields.schedule, gotTask.Schedule())
}
