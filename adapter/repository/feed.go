package repository

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/inhies/go-bytesize"

	"go-feedmaker/entity"
	"go-feedmaker/interactor"
)

type (
	FeedConfig struct {
		CountQuery    string
		SelectQuery   string
		FileSizeLimit bytesize.ByteSize
		FileLineLimit uint
		SqlGateway    SqlGateway
	}

	RedisClient interface {
		Connection() Connection
		PubSub() PubSub
	}

	Connection interface {
		redis.Conn
	}

	PubSub interface {
		Subscribe(channel ...interface{}) error
		Unsubscribe(channel ...interface{}) error
		Receive() interface{}
		Ping(data string) error
		io.Closer
	}

	feedRepo struct {
		client         RedisClient
		idSetName      string
		cancelChanName string
		typeConfigMap  map[string]*FeedConfig
		ftpGateway     FtpGateway
	}
)

func NewFeedRepo(config map[string]*FeedConfig, client RedisClient, ftpGateway FtpGateway) *feedRepo {
	return &feedRepo{
		client:        client,
		idSetName:     "generationIDs",
		ftpGateway:    ftpGateway,
		typeConfigMap: config,
	}
}

func (r *feedRepo) GetFactoryByGenerationType(generationType string) (interactor.FeedFactory, error) {
	config, ok := r.typeConfigMap[generationType]
	if !ok {
		return nil, entity.ErrInvalidGenerationType
	}

	return NewDefaultFactory(config, config.SqlGateway, r.ftpGateway, generationType)
}

func (r *feedRepo) ListAllowedTypes() []string {
	types := make([]string, 0, len(r.typeConfigMap))
	for k := range r.typeConfigMap {
		types = append(types, k)
	}
	return types
}

func (r *feedRepo) IsAllowedType(generationType string) bool {
	_, ok := r.typeConfigMap[generationType]
	return ok
}

func (r *feedRepo) StoreGeneration(ctx context.Context, generation *entity.Generation) error {
	conn := r.client.Connection()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("SADD", r.idSetName, generation.ID)
	hashArgs := new(redis.Args).
		Add(generation.ID).
		Add("type", generation.Type).
		Add("progress", generation.Progress).
		Add("data_fetched", generation.DataFetched).
		Add("files_uploaded", generation.FilesUploaded).
		Add("start_time", generation.StartTime.Unix())
	if !generation.EndTime.IsZero() {
		hashArgs = hashArgs.Add("end_time", generation.EndTime.Unix())
	}
	conn.Send("HMSET", hashArgs...)

	_, err := conn.Do("EXEC")
	return err
}

func (r *feedRepo) GetGeneration(ctx context.Context, generationID string) (*entity.Generation, error) {
	conn := r.client.Connection()
	defer conn.Close()
	stringMap, err := redis.StringMap(conn.Do("HGETALL", generationID))
	if err != nil {
		return nil, err
	}
	stringMap["id"] = generationID
	generation, err := makeGenerationFromRedisValues(stringMap)
	if err != nil {
		return nil, err
	}
	return generation, nil
}

func (r *feedRepo) ListGenerations(ctx context.Context) ([]*entity.Generation, error) {
	generations := make([]*entity.Generation, 0)
	conn := r.client.Connection()
	defer conn.Close()
	generationIDs, err := redis.Strings(conn.Do("SMEMBERS", r.idSetName))
	if err != nil {
		return nil, err
	}
	for _, id := range generationIDs {
		generation, err := r.GetGeneration(ctx, id)
		if err != nil {
			return nil, err
		}
		generations = append(generations, generation)
	}
	return generations, nil
}

func makeGenerationFromRedisValues(v map[string]string) (*entity.Generation, error) {
	progress, _ := strconv.ParseUint(v["progress"], 10, 32)
	filesUploaded, _ := strconv.ParseUint(v["files_uploaded"], 10, 32)
	dataFetched, _ := strconv.ParseBool(v["data_fetched"])
	isCanceled, _ := strconv.ParseBool(v["is_canceled"])

	generation := new(entity.Generation)
	generation.ID = v["id"]
	generation.Type = v["type"]
	generation.Progress = uint(progress)
	generation.DataFetched = dataFetched
	generation.IsCanceled = isCanceled
	generation.FilesUploaded = uint(filesUploaded)

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

func (r *feedRepo) UpdateGenerationState(ctx context.Context, generation *entity.Generation) error {
	channel := "generation.updated"
	conn := r.client.Connection()
	defer conn.Close()
	hashArgs := new(redis.Args).Add(generation.ID).
		Add("progress", generation.Progress).
		Add("data_fetched", generation.DataFetched).
		Add("files_uploaded", generation.FilesUploaded)
	if !generation.EndTime.IsZero() {
		hashArgs = hashArgs.Add("end_time", generation.EndTime.Unix())
	}
	_, err := conn.Do("HSET", hashArgs...)
	if err != nil {
		return err
	}
	_, err = conn.Do("PUBLISH", channel, generation.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *feedRepo) DeleteGeneration(ctx context.Context, generationID string) error {
	conn := r.client.Connection()
	defer conn.Close()

	_, err := conn.Do("DEL", generationID)
	return err
}

func (r *feedRepo) CancelGeneration(ctx context.Context, id string) error {
	conn := r.client.Connection()
	defer conn.Close()
	if _, err := conn.Do("HSET", id, "is_canceled", true); err != nil {
		return err
	}
	channel := fmt.Sprintf("%s.canceled", id)
	_, err := conn.Do("PUBLISH", channel, "1")
	return err
}

func (r *feedRepo) OnGenerationCanceled(ctx context.Context, generationID string, callback func()) error {
	channel := fmt.Sprintf("%s.canceled", generationID)
	pubsub := r.client.PubSub()
	defer pubsub.Close()
	if err := pubsub.Subscribe(channel); err != nil {
		return err
	}
	defer pubsub.Unsubscribe(channel)
	errChan := make(chan error)

	go func() {
		for {
			switch v := pubsub.Receive().(type) {
			case error:
				errChan <- v
				return
			case redis.Message:
				if v.Channel == channel {
					callback()
					errChan <- nil
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := pubsub.Ping(""); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			return err
		}
	}
}

func (r *feedRepo) OnGenerationsUpdated(ctx context.Context, callback func(*entity.Generation)) error {
	channel := "generation.updated"
	errChan := make(chan error)
	pubsub := r.client.PubSub()
	defer pubsub.Close()

	if err := pubsub.Subscribe(channel); err != nil {
		return err
	}
	defer pubsub.Unsubscribe(channel)

	go func() {
		for {
			switch v := pubsub.Receive().(type) {
			case error:
				errChan <- v
				return
			case redis.Message:
				if v.Channel == channel {
					if generation, err := r.GetGeneration(ctx, string(v.Data)); err != nil {
						errChan <- err
					} else {
						callback(generation)
					}
				}
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := pubsub.Ping(""); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			return err
		}
	}
}
