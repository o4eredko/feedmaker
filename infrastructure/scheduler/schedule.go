package scheduler

import (
	"time"

	"github.com/rs/zerolog/log"
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
		startTimestamp := s.getAlignedStartTimestamp(nowTimestamp)
		log.Info().
			Str("now", nowTimestamp.UTC().Format(time.RFC3339)).
			Str("timestamp", startTimestamp.UTC().Format(time.RFC3339)).
			Msg("first call of Next")
		return startTimestamp
	}

	nextTimestamp := s.getNextTimestamp(nowTimestamp)
	log.Info().
		Str("now", nowTimestamp.UTC().Format(time.RFC3339)).
		Str("timestamp", nextTimestamp.UTC().Format(time.RFC3339)).
		Msg("another one call of Next")
	return nextTimestamp
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
