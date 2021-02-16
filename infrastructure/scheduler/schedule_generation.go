package scheduler

import (
	"context"

	"go-feedmaker/interactor"
)

func (s *Scheduler) ScheduleGeneration(feeds interactor.FeedInteractor, generationType string, schedule *Schedule) error {
	cmd, err := NewCmd(feeds.GenerateFeed, context.Background(), generationType)
	if err != nil {
		return err
	}
	task := NewTask(cmd, schedule)
	return s.ScheduleTask(TaskID(generationType), task)
}

func (s *Scheduler) ScheduleAllSavedGenerations(feeds interactor.FeedInteractor) error {
	schedules, err := s.ListSchedules()
	if err != nil {
		return err
	}
	for taskID, schedule := range schedules {
		if err := s.ScheduleGeneration(feeds, string(taskID), schedule); err != nil {
			return err
		}
	}
	return nil
}
