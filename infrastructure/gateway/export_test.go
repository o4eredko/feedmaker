package gateway

import "github.com/gomodule/redigo/redis"

func (f *FtpGateway) SetConnection(connection FtpConnection) {
	f.connection = connection
}

func (r *RedisGateway) SetPool(pool *redis.Pool) {
	r.pool = pool
}
