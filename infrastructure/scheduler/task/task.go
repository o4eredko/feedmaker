package task

import "go-feedmaker/infrastructure/scheduler"

type (
	Task struct {
		cmd      scheduler.Runner
		schedule scheduler.Nexter
	}
)

func New(cmd scheduler.Runner, schedule scheduler.Nexter) *Task {
	return &Task{
		cmd:      cmd,
		schedule: schedule,
	}
}

func (t *Task) Cmd() scheduler.Runner {
	return t.cmd
}

func (t *Task) Schedule() scheduler.Nexter {
	return t.schedule
}
