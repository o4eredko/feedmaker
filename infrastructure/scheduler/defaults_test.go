package scheduler_test

import (
	"errors"
	"time"

	"github.com/robfig/cron/v3"

	"go-feedmaker/infrastructure/scheduler"
)

var (
	defaultErr              = errors.New("default error")
	defaultEntryID          = cron.EntryID(42)
	defaultTaskID           = scheduler.TaskID("foobar")
	defaultSchedule         = scheduler.NewSchedule(time.Now().UTC().Truncate(time.Second), time.Second*42)
	defaultScheduledTaskIDs = []scheduler.TaskID{defaultTaskID, "spam", "ham", "eggs"}
	defaultTaskSchedules    = makeTaskSchedules(defaultScheduledTaskIDs)
)

func makeTaskSchedules(ids []scheduler.TaskID) map[scheduler.TaskID]*scheduler.Schedule {
	schedules := make(map[scheduler.TaskID]*scheduler.Schedule, len(ids))
	for i, id := range ids {
		duration := time.Duration(i)
		start := time.Now().Add(4 + time.Hour*duration)
		interval := 2 + time.Hour*duration
		schedules[id] = scheduler.NewSchedule(start, interval)
	}
	return schedules
}
