package entity

import "time"

type Generation struct {
	ID        string
	Type      string
	Progress  uint
	StartTime time.Time
	EndTime   time.Time
}

func (g *Generation) SetProgress(progress uint) {
	if progress > 100 {
		g.Progress = 100
	} else {
		g.Progress = progress
	}
}
