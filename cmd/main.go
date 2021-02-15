package main

import (
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"

	"go-feedmaker/infrastructure/gateway"
)

type redisDialer struct{}

func (r *redisDialer) Dial(network, addr string, options ...redis.DialOption) (gateway.RedisConnection, error) {
	return redis.Dial(network, addr, options...)
}

func main() {
	log.Info().Msgf("Hello world")
}
