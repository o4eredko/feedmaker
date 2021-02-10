package main

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"

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
	defer redisGateway.Disconnect()

	repo := repository.NewFeedRepo(redisGateway)

	generation := &entity.Generation{
		ID:        "abc",
		Type:      "test type",
		Progress:  50,
		StartTime: time.Now(),
	}

	// go func() {
	// 	time.Sleep(time.Second * 3)
	// 	if err := repo.CancelGeneration(context.Background(), generation.ID); err != nil {
	// 		panic(err)
	// 	}
	// }()
	callback := func() {
		log.Info().Msgf("Generation %s was canceled", generation.ID)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := repo.OnGenerationCanceled(ctx, generation.ID, callback); err != nil {
		panic(err)
	}
	log.Info().Msgf("returned")
}
