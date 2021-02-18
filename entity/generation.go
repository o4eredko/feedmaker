package entity

import "time"

type Generation struct {
	ID            string
	Type          string
	Progress      uint
	DataFetched   bool
	FilesUploaded uint
	IsCanceled    bool
	StartTime     time.Time
	EndTime       time.Time
}

func (g *Generation) SetProgress(progress uint) {
	if progress > 100 {
		progress = 100
	}
	g.Progress = progress
	if g.Progress == 100 {
		g.EndTime = time.Now()
	}
}
