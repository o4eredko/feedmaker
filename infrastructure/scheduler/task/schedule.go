package task

import "time"

type (
	ScheduleWithSpecificStartTimestamp struct {
		StartTimestamp         time.Time
		DelayInterval          time.Duration
		startTimestampExceeded bool
	}
)

func NewScheduleWithSpecificStartTimestamp(
	startTimestamp time.Time,
	delayInterval time.Duration,
) *ScheduleWithSpecificStartTimestamp {
	return &ScheduleWithSpecificStartTimestamp{
		StartTimestamp: startTimestamp,
		DelayInterval:  delayInterval,
	}
}

func (s *ScheduleWithSpecificStartTimestamp) Next(nowTimestamp time.Time) time.Time {
	if !s.startTimestampExceeded {
		s.startTimestampExceeded = true
		return s.getAlignedStartTimestamp(nowTimestamp)
	}

	return s.getNextTimestamp(nowTimestamp)
}

func (s *ScheduleWithSpecificStartTimestamp) getAlignedStartTimestamp(nowTimestamp time.Time) time.Time {
	startTimestamp := s.StartTimestamp
	if startTimestamp.After(nowTimestamp) {
		return startTimestamp
	}
	elapsedTimeRoundedToInterval := nowTimestamp.
		Sub(startTimestamp).
		Round(s.DelayInterval)
	return startTimestamp.Add(elapsedTimeRoundedToInterval)
}

func (s *ScheduleWithSpecificStartTimestamp) getNextTimestamp(nowTimestamp time.Time) time.Time {
	return nowTimestamp.
		Add(s.DelayInterval).
		Truncate(time.Second)
}
