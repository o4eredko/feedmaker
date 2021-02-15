package scheduler

type (
	Task struct {
		Cmd      Runner
		Schedule *Schedule
	}

	Runner interface {
		Run()
	}
)

func NewTask(cmd Runner, schedule *Schedule) *Task {
	return &Task{
		Cmd:      cmd,
		Schedule: schedule,
	}
}
