package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"

	"go-feedmaker/entity"
	"go-feedmaker/interactor"
)

type (
	RedisClient interface {
		Do(commandName string, args ...interface{}) (reply interface{}, err error)
		Send(commandName string, args ...interface{}) error
		Flush() error
		Receive() (reply interface{}, err error)
	}

	feedRepo struct {
		client    RedisClient
		idSetName string
	}
)

func NewFeedRepo(client RedisClient) *feedRepo {
	return &feedRepo{
		client:    client,
		idSetName: "generationIDs",
	}
}

func (r *feedRepo) GetFactoryByGenerationType(generationType string) (interactor.FeedFactory, error) {
	switch {
	case strings.HasPrefix(generationType, "yandex-"):
		return NewYandexFactory(), nil
	default:
		return NewDefaultFactory(), nil
	}
}

func (r *feedRepo) StoreGeneration(ctx context.Context, generation *entity.Generation) error {
	r.client.Send("MULTI")
	r.client.Send("SADD", r.idSetName, generation.ID)
	hashArgs := new(redis.Args).
		Add(generation.ID).
		Add("type", generation.Type).
		Add("progress", generation.Progress).
		Add("start_time", generation.StartTime.Unix())
	if !generation.EndTime.IsZero() {
		hashArgs = hashArgs.Add("end_time", generation.EndTime.Unix())
	}
	r.client.Send("HMSET", hashArgs...)

	_, err := r.client.Do("EXEC")
	return err
}

func (r *feedRepo) ListGenerations(ctx context.Context) ([]*entity.Generation, error) {
	generations := make([]*entity.Generation, 0)
	generationIDs, err := redis.Strings(r.client.Do("SMEMBERS", r.idSetName))
	if err != nil {
		return nil, err
	}
	for _, id := range generationIDs {
		stringMap, err := redis.StringMap(r.client.Do("HGETALL", id))
		if err != nil {
			return nil, err
		}
		stringMap["id"] = id
		generation, err := makeGenerationFromRedisValues(stringMap)
		if err != nil {
			return nil, err
		}
		generations = append(generations, generation)
	}
	return generations, nil
}

func makeGenerationFromRedisValues(v map[string]string) (*entity.Generation, error) {
	generation := new(entity.Generation)
	generation.ID = v["id"]
	generation.Type = v["type"]
	if timestamp, ok := v["start_time"]; ok && len(timestamp) > 0 {
		startTime, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%s 'start_time': %w", generation.ID, entity.ErrInvalidTimestamp)
		}
		generation.StartTime = time.Unix(startTime, 0)
	}
	if timestamp, ok := v["end_time"]; ok && len(timestamp) > 0 {
		startTime, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%s 'end_time': %w", generation.ID, entity.ErrInvalidTimestamp)
		}
		generation.EndTime = time.Unix(startTime, 0)
	}
	return generation, nil
}

func (r *feedRepo) UpdateProgress(ctx context.Context, generation *entity.Generation) error {
	hashArgs := new(redis.Args).Add(generation.ID).Add("progress", generation.Progress)
	if !generation.EndTime.IsZero() {
		hashArgs = hashArgs.Add("end_time", generation.EndTime.Unix())
	}
	_, err := r.client.Do("HSET", hashArgs...)
	if err != nil {
		return err
	}
	_, err = r.client.Do("PUBLISH", generation.ID, generation.Progress)
	if err != nil {
		return err
	}
	return nil
}

func (r *feedRepo) ListAllowedTypes(ctx context.Context) ([]string, error) {
	panic("implement me")
}

func (r *feedRepo) IsAllowedType(ctx context.Context, generationType string) (bool, error) {
	panic("implement me")
}
