package entity

import "time"

type Generation struct {
	ID          string
	Type        string
	IsCompleted bool
	StartTime   time.Time
	EndTime     time.Time
}
