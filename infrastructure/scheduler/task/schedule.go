package task

import "time"

type (
	ScheduleWithSpecificStartTimestamp struct {
		startTimestamp         time.Time
		fireInterval           time.Duration
		startTimestampExceeded bool
	}
)

func NewScheduleWithSpecificStartTimestamp(
	startTimestamp time.Time,
	delayInterval time.Duration,
) *ScheduleWithSpecificStartTimestamp {
	return &ScheduleWithSpecificStartTimestamp{
		startTimestamp: startTimestamp,
		fireInterval:   delayInterval,
	}
}

func (s *ScheduleWithSpecificStartTimestamp) StartTimestamp() time.Time {
	return s.startTimestamp
}

func (s *ScheduleWithSpecificStartTimestamp) FireInterval() time.Duration {
	return s.fireInterval
}

func (s *ScheduleWithSpecificStartTimestamp) Next(nowTimestamp time.Time) time.Time {
	if !s.startTimestampExceeded {
		s.startTimestampExceeded = true
		return s.getAlignedStartTimestamp(nowTimestamp)
	}

	return s.getNextTimestamp(nowTimestamp)
}

func (s *ScheduleWithSpecificStartTimestamp) getAlignedStartTimestamp(nowTimestamp time.Time) time.Time {
	startTimestamp := s.startTimestamp
	if startTimestamp.After(nowTimestamp) {
		return startTimestamp
	}
	elapsedTimeRoundedToInterval := nowTimestamp.
		Sub(startTimestamp).
		Round(s.fireInterval)
	return startTimestamp.Add(elapsedTimeRoundedToInterval)
}

func (s *ScheduleWithSpecificStartTimestamp) getNextTimestamp(nowTimestamp time.Time) time.Time {
	return nowTimestamp.
		Add(s.fireInterval).
		Truncate(time.Second)
}
