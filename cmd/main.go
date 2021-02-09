package main

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"

	"go-feedmaker/adapter/repository"
	"go-feedmaker/entity"
	"go-feedmaker/infrastructure/gateway"
)

type redisDialer struct{}

func (r *redisDialer) Dial(network, addr string, options ...redis.DialOption) (gateway.RedisConnection, error) {
	return redis.Dial(network, addr, options...)
}

func main() {
	redisGateway := &gateway.RedisGateway{
		Config: gateway.RedisConfig{
			Host:        "localhost",
			Port:        "6379",
			ConnTimeout: time.Second,
		},
		Dialer: new(redisDialer),
	}
	if err := redisGateway.Connect(); err != nil {
		panic(err)
	}

	repo := repository.NewFeedRepo(redisGateway)
	generation := &entity.Generation{
		ID:        uuid.New().String(),
		Type:      "test type",
		Progress:  100,
		StartTime: time.Now(),
		EndTime:   time.Now(),
	}

	if err := repo.StoreGeneration(context.Background(), generation); err != nil {
		panic(err)
	}

	if err := redisGateway.Disconnect(); err != nil {
		panic(err)
	}
}
