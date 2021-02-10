package scheduler_test

import (
	"errors"

	"github.com/robfig/cron/v3"

	"go-feedmaker/infrastructure/scheduler"
)

var (
	defaultErr     = errors.New("default error")
	defaultEntryID = cron.EntryID(42)
	defaultTaskID  = scheduler.TaskID("foobar")
)
