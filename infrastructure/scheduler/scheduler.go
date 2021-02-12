package scheduler

import (
	"context"

	"github.com/robfig/cron/v3"
)

type (
	Scheduler struct {
		cron          Croner
		scheduleSaver ScheduleSaver
		mapper        TaskIDMapper
	}

	Croner interface {
		Start()
		Stop() context.Context
		Schedule(cron.Schedule, cron.Job) cron.EntryID
		Remove(cron.EntryID)
	}

	TaskID string

	ScheduleSaver interface {
		Store(TaskID, *Schedule) error
		Load(TaskID) (*Schedule, error)
		Delete(TaskID) error
		ListScheduledTaskIDs() ([]TaskID, error)
	}

	TaskIDMapper interface {
		Store(TaskID, cron.EntryID) error
		Load(TaskID) (cron.EntryID, error)
		Delete(TaskID) error
	}
)

func New(cron Croner, scheduleSaver ScheduleSaver) *Scheduler {
	return &Scheduler{
		cron:          cron,
		scheduleSaver: scheduleSaver,
		mapper:        NewMapper(),
	}
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

func (s *Scheduler) ScheduleTask(taskID TaskID, task *Task) error {
	if err := s.scheduleSaver.Store(taskID, task.Schedule); err != nil {
		return err
	}
	entryID := s.cron.Schedule(task.Schedule, task.Cmd)
	if err := s.mapper.Store(taskID, entryID); err != nil {
		s.cron.Remove(entryID)
		return err
	}
	return nil
}

func (s *Scheduler) RemoveTask(taskID TaskID) error {
	entryID, err := s.mapper.Load(taskID)
	if err != nil {
		return err
	}
	s.cron.Remove(entryID)
	if err := s.mapper.Delete(taskID); err != nil {
		return err
	}
	return s.scheduleSaver.Delete(taskID)
}

func (s *Scheduler) ListSchedules() (map[TaskID]*Schedule, error) {
	ids, err := s.scheduleSaver.ListScheduledTaskIDs()
	if err != nil {
		return nil, err
	}
	schedules := make(map[TaskID]*Schedule)
	for _, id := range ids {
		schedule, err := s.scheduleSaver.Load(id)
		if err != nil {
			return nil, err
		}
		schedules[id] = schedule
	}
	return schedules, err
}
