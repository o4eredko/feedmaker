package rest

import "time"

type (
	scheduleTaskIn struct {
		StartTimestamp time.Time     `json:"start_timestamp"`
		DelayInterval  time.Duration `json:"delay_interval"`
	}

	scheduleOut struct {
		GenerationType string        `json:"generation_type"`
		StartTimestamp time.Time     `json:"start_timestamp"`
		DelayInterval  time.Duration `json:"delay_interval"`
	}
)
