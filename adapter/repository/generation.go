package repository

import (
	"context"
	"strings"

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
		hashName  string
	}
)

func NewFeedRepo(client RedisClient) *feedRepo {
	return &feedRepo{
		client:    client,
		idSetName: "generationIDs",
		hashName:  "generations",
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
	hashArgs := new(redis.Args).Add(generation.ID).AddFlat(generation)
	r.client.Send("HMSET", hashArgs...)

	_, err := r.client.Do("EXEC")
	return err
}

func (r *feedRepo) ListGenerations(ctx context.Context) ([]*entity.Generation, error) {
	// res := make([]*entity.Generation, 0)
	panic("implement me")
}

func (r *feedRepo) ListAllowedTypes(ctx context.Context) ([]string, error) {
	panic("implement me")
}

func (r *feedRepo) IsAllowedType(ctx context.Context, generationType string) (bool, error) {
	panic("implement me")
}
