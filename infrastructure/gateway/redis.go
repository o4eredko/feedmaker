package gateway

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/mediocregopher/radix/v3"
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
		Dial(network, addr string, options ...radix.DialOpt) (RedisConnection, error)
	}

	RedisConnection interface {
		radix.Conn
	}

	RedisClient interface {
		radix.Client
	}

	RedisGateway struct {
		Config RedisConfig
		Dialer RedisDialer
		pool   RedisClient
	}
)

func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (r *RedisGateway) connect(network, addr string) (radix.Conn, error) {
	return r.Dialer.Dial(network, addr, radix.DialConnectTimeout(r.Config.ConnTimeout))
}

func (r *RedisGateway) Connect() error {
	pool, err := radix.NewPool(
		"tcp", r.Config.Addr(),
		r.Config.PoolSize,
		radix.PoolConnFunc(r.connect),
	)
	if err != nil {
		return err
	}
	r.pool = pool
	return nil
}

func (r *RedisGateway) Do(action radix.Action) error {
	if reflect.ValueOf(r.pool).IsNil() {
		return ErrRedisDisconnected
	}
	return r.pool.Do(action)
}

func (r *RedisGateway) Disconnect() error {
	if reflect.ValueOf(r.pool).IsNil() {
		return ErrRedisDisconnected
	}
	return r.pool.Close()
}
