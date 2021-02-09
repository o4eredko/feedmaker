package gateway

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	ErrRedisDisconnected = errors.New("gateway is not connected to Redis")
)

type (
	RedisConfig struct {
		Host        string
		Port        string
		ConnTimeout time.Duration
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

	RedisConnection interface {
		redis.Conn
	}

	RedisGateway struct {
		Config     RedisConfig
		Dialer     RedisDialer
		connection RedisConnection
		RedisClient
	}
)

func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (r *RedisGateway) Connect() error {
	conn, err := r.Dialer.Dial("tcp", r.Config.Addr(), redis.DialConnectTimeout(r.Config.ConnTimeout))
	if err != nil {
		return err
	}
	r.connection = conn
	r.RedisClient = conn
	return nil
}

func (r *RedisGateway) Disconnect() error {
	if reflect.ValueOf(r.connection).IsNil() {
		return ErrRedisDisconnected
	}
	return r.connection.Close()
}
