package scheduler

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
)

type (
	TaskScheduler interface {
		StartStopper
		AddTask(taskID TaskID, task Task) error
		RemoveTask(taskID TaskID) error
	}

	StartStopper interface {
		Start()
		Stop() context.Context
	}

	TaskID string

	Task interface {
		Cmd() Runner
		Schedule() Schedule
	}

	Runner interface {
		Run()
	}

	Schedule interface {
		Nexter
		StartTimestamp() time.Time
		FireInterval() time.Duration
	}

	Nexter interface {
		Next(time.Time) time.Time
	}

	Scheduler struct {
		cron   Croner
		mapper TaskIDMapper
	}

	TaskIDMapper interface {
		Store(TaskID, cron.EntryID) error
		Load(TaskID) (cron.EntryID, error)
		Delete(TaskID) error
	}

	Croner interface {
		StartStopper
		Schedule(cron.Schedule, cron.Job) cron.EntryID
		Remove(cron.EntryID)
	}
)

func New(cron Croner, mapper TaskIDMapper) *Scheduler {
	return &Scheduler{
		mapper: mapper,
		cron:   cron,
	}
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() context.Context {
	return s.cron.Stop()
}

func (s *Scheduler) AddTask(taskID TaskID, task Task) error {
	if _, err := s.mapper.Load(taskID); err != nil {
		return err
	}
	entryID := s.cron.Schedule(task.Schedule(), task.Cmd())
	return s.mapper.Store(taskID, entryID)
}

func (s *Scheduler) RemoveTask(taskID TaskID) error {
	entryID, err := s.mapper.Load(taskID)
	if err != nil {
		return err
	}
	s.cron.Remove(entryID)
	return s.mapper.Delete(taskID)
}
