package scheduler

import (
	"errors"
	"sync"

	"github.com/robfig/cron/v3"
)

type (
	Mapper struct {
		readWriteLocker sync.RWMutex
		mapping         map[TaskID]cron.EntryID
	}
)

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrTaskAlreadyExists = errors.New("task already exists")
)

func NewMapper() *Mapper {
	return &Mapper{
		mapping: make(map[TaskID]cron.EntryID),
	}
}

func (m *Mapper) Store(taskID TaskID, entryID cron.EntryID) error {
	m.readWriteLocker.Lock()
	defer m.readWriteLocker.Unlock()
	if _, exists := m.mapping[taskID]; exists {
		return ErrTaskAlreadyExists
	}
	m.mapping[taskID] = entryID
	return nil
}

func (m *Mapper) Load(taskID TaskID) (cron.EntryID, error) {
	m.readWriteLocker.RLock()
	defer m.readWriteLocker.RUnlock()
	entryID, exists := m.mapping[taskID]
	if !exists {
		return 0, ErrTaskNotFound
	}
	return entryID, nil
}

func (m *Mapper) Delete(taskID TaskID) error {
	m.readWriteLocker.Lock()
	defer m.readWriteLocker.Unlock()
	if _, exists := m.mapping[taskID]; !exists {
		return ErrTaskNotFound
	}
	delete(m.mapping, taskID)
	return nil
}
