package scheduler

import "github.com/robfig/cron/v3"

func (s *Scheduler) Cron() Croner {
	return s.cron
}

func (s *Scheduler) TaskIDMapping() TaskIDMapper {
	return s.taskIDMapping
}

func (m *Mapper) SetMapping(mapping map[TaskID]cron.EntryID) {
	m.taskIDMapping = mapping
}
