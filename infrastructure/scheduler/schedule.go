package scheduler

import (
	"time"
)

type (
	Schedule struct {
		startTimestamp         time.Time
		fireInterval           time.Duration
		startTimestampExceeded bool
	}
)

func NewSchedule(
	startTimestamp time.Time,
	delayInterval time.Duration,
) *Schedule {
	return &Schedule{
		startTimestamp: startTimestamp,
		fireInterval:   delayInterval,
	}
}

func (s *Schedule) StartTimestamp() time.Time {
	return s.startTimestamp
}

func (s *Schedule) FireInterval() time.Duration {
	return s.fireInterval
}

func (s *Schedule) Next(nowTimestamp time.Time) time.Time {
	if !s.startTimestampExceeded {
		s.startTimestampExceeded = true
		return s.getAlignedStartTimestamp(nowTimestamp)
	}
	return s.getNextTimestamp(nowTimestamp)
}

func (s *Schedule) getAlignedStartTimestamp(nowTimestamp time.Time) time.Time {
	startTimestamp := s.startTimestamp
	if startTimestamp.After(nowTimestamp) {
		return startTimestamp
	}
	elapsedTimeRoundedToInterval := nowTimestamp.Sub(startTimestamp).Round(s.fireInterval)
	startTimestamp = startTimestamp.Add(elapsedTimeRoundedToInterval)
	if startTimestamp.Before(nowTimestamp) {
		startTimestamp = startTimestamp.Add(s.fireInterval)
	}
	return startTimestamp
}

func (s *Schedule) getNextTimestamp(nowTimestamp time.Time) time.Time {
	return nowTimestamp.
		Add(s.fireInterval).
		Truncate(time.Second)
}
