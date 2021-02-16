package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jlaffaye/ftp"
	"github.com/rs/zerolog/log"

	"go-feedmaker/infrastructure/config"
	"go-feedmaker/infrastructure/gateway"
)

type (
	redisDialer struct{}
	ftpDialer   struct{}
)

func (r *redisDialer) Dial(network, addr string, options ...redis.DialOption) (gateway.RedisConnection, error) {
	return redis.Dial(network, addr, options...)
}

func (f *ftpDialer) DialTimeout(addr string, timeout time.Duration) (gateway.FtpConnection, error) {
	return ftp.DialTimeout(addr, timeout)
}

func main() {
	conf := config.LoadFromFile(config.DefaultConfig)
	if err := conf.Logger.Apply(); err != nil {
		log.Fatal().Err(err).Msgf("Cannot configure logger")
	}

	redisGateway := &gateway.RedisGateway{
		Config: conf.Redis,
		Dialer: new(redisDialer),
	}
	if err := redisGateway.Connect(); err != nil {
		log.Fatal().Err(err).Msgf("Can't connect to Redis")
	}
	defer redisGateway.Disconnect()

	ftpGateway := &gateway.FtpGateway{
		Dialer: new(ftpDialer),
		Config: conf.Ftp,
	}
	if err := ftpGateway.Connect(); err != nil {
		log.Fatal().Err(err).Msgf("Can't connect to Ftp")
	}
	defer ftpGateway.Disconnect()

	log.Info().Msgf("Successfully started server")
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info().Msgf("Server was stopped")
}
