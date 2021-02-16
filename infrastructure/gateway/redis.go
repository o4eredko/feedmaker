package gateway

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"

	"go-feedmaker/adapter/repository"
)

var (
	ErrRedisDisconnected = errors.New("gateway is not connected to Redis")
)

type (
	RedisConfig struct {
		Host        string
		Port        string
		ConnTimeout time.Duration
		PoolSize    int
	}

	RedisDialer interface {
		Dial(network, addr string, options ...redis.DialOption) (RedisConnection, error)
	}

	RedisClient interface {
		Do(commandName string, args ...interface{}) (reply interface{}, err error)
		Send(commandName string, args ...interface{}) error
		Flush() error
		Receive() (reply interface{}, err error)
	}

	PubSub interface {
		Subscribe(channel ...interface{}) error
		Unsubscribe(channel ...interface{}) error
		Receive() interface{}
	}

	RedisConnection interface {
		redis.Conn
	}

	RedisGateway struct {
		Config     RedisConfig
		Dialer     RedisDialer
		connection RedisConnection
		pool       *redis.Pool
	}
)

func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (r *RedisGateway) dial() (redis.Conn, error) {
	conn, err := r.Dialer.Dial("tcp", r.Config.Addr(), redis.DialConnectTimeout(r.Config.ConnTimeout))
	return conn, err
}

func (r *RedisGateway) Connect() error {
	r.pool = &redis.Pool{
		Dial:        r.dial,
		MaxIdle:     r.Config.PoolSize,
		MaxActive:   r.Config.PoolSize,
		IdleTimeout: time.Minute,
		Wait:        true,
	}
	return r.ping()
}

func (r *RedisGateway) ping() error {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	return err
}

func (r *RedisGateway) Connection() repository.Connection {
	return r.pool.Get()
}

func (r *RedisGateway) PubSub() repository.PubSub {
	return &redis.PubSubConn{Conn: r.pool.Get()}
}

func (r *RedisGateway) Disconnect() error {
	if r.pool == nil {
		return ErrRedisDisconnected
	}
	return r.pool.Close()
}
