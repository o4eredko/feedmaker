package scheduler

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"

	"go-feedmaker/adapter/repository"
)

type (
	scheduleSaver struct {
		client repository.RedisClient
	}
)

const (
	TaskIDsKey = "task_ids"
)

var (
	ErrInvalidTimestamp = errors.New("invalid timestamp")
	ErrInvalidInterval  = errors.New("invalid interval")
)

func NewScheduleSaver(client repository.RedisClient) *scheduleSaver {
	return &scheduleSaver{
		client: client,
	}
}

func (s *scheduleSaver) Store(id TaskID, schedule *Schedule) error {
	conn := s.client.Connection()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("SADD", TaskIDsKey, id)
	args := makeRedisArgs(id, schedule)
	conn.Send("HMSET", args...)
	_, err := conn.Do("EXEC")
	return err
}

func (s *scheduleSaver) Load(id TaskID) (*Schedule, error) {
	conn := s.client.Connection()
	defer conn.Close()
	rawSchedule, err := redis.StringMap(conn.Do("HGETALL", id))
	if err != nil {
		return nil, err
	}
	return makeSchedule(rawSchedule)
}

func (s scheduleSaver) Delete(id TaskID) error {
	conn := s.client.Connection()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("SREM", TaskIDsKey, id)
	conn.Send("DEL", id)
	_, err := conn.Do("EXEC")
	return err
}

func (s scheduleSaver) ListScheduledTaskIDs() ([]TaskID, error) {
	conn := s.client.Connection()
	defer conn.Close()
	taskIDs := make([]TaskID, 0)
	rawTaskIDs, err := redis.Strings(conn.Do("SMEMBERS", TaskIDsKey))
	if err != nil {
		return nil, err
	}
	for _, id := range rawTaskIDs {
		taskIDs = append(taskIDs, TaskID(id))
	}
	return taskIDs, nil
}

func makeRedisArgs(id TaskID, schedule *Schedule) redis.Args {
	return new(redis.Args).
		Add(id).
		Add("start_timestamp", schedule.StartTimestamp().Unix()).
		Add("fire_interval", schedule.FireInterval().Seconds())
}

func makeSchedule(v map[string]string) (*Schedule, error) {
	schedule := new(Schedule)
	rawStartTimestamp := v["start_timestamp"]
	startTimestamp, err := strconv.ParseInt(rawStartTimestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("start_timestamp: %v: %w",
			rawStartTimestamp, ErrInvalidTimestamp)
	}
	schedule.startTimestamp = time.Unix(startTimestamp, 0).UTC()
	rawFireInterval := v["fire_interval"]
	fireInterval, err := strconv.ParseInt(rawFireInterval, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("fire_interval: %v: %w",
			rawFireInterval, ErrInvalidInterval)
	}
	schedule.fireInterval = time.Second * time.Duration(fireInterval)
	return schedule, nil
}
