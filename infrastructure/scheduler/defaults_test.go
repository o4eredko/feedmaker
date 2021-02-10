package scheduler_test

import (
	"errors"
	"time"

	"github.com/robfig/cron/v3"

	"go-feedmaker/infrastructure/scheduler"
	"go-feedmaker/infrastructure/scheduler/task"
)

var (
	defaultErr              = errors.New("default error")
	defaultEntryID          = cron.EntryID(42)
	defaultTaskID           = scheduler.TaskID("foobar")
	defaultSchedule         = task.NewSchedule(time.Now().UTC(), time.Second*42)
	defaultScheduledTaskIDs = []scheduler.TaskID{defaultTaskID, "spam", "ham", "eggs"}
	defaultTaskSchedules    = makeTaskSchedules(defaultScheduledTaskIDs)
)

func makeTaskSchedules(ids []scheduler.TaskID) map[scheduler.TaskID]*task.Schedule {
	schedules := make(map[scheduler.TaskID]*task.Schedule, len(ids))
	for i, id := range ids {
		duration := time.Duration(i)
		start := time.Now().Add(4 + time.Hour*duration)
		interval := 2 + time.Hour*duration
		schedules[id] = task.NewSchedule(start, interval)
	}
	return schedules
}
