package rest

import "time"

type (
	scheduleTaskIn struct {
		StartTimestamp time.Time `json:"start_timestamp"`
		DelayInterval  int       `json:"delay_interval"`
	}

	scheduleOut struct {
		StartTimestamp string `json:"start_timestamp"`
		DelayInterval  int    `json:"delay_interval"`
	}
)
