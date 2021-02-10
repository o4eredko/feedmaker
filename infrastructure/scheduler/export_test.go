package scheduler

import "github.com/robfig/cron/v3"

func (s *Scheduler) Cron() Croner {
	return s.cron
}

func (s *Scheduler) Saver() ScheduleSaver {
	return s.scheduleSaver
}

func (s *Scheduler) Mapper() TaskIDMapper {
	return s.mapper
}

func (s *Scheduler) SetMapper(mapper TaskIDMapper) {
	s.mapper = mapper
}

func (m *Mapper) SetMapping(mapping map[TaskID]cron.EntryID) {
	m.mapping = mapping
}
