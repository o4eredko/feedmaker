package gateway

func (f *FtpGateway) SetConnection(connection FtpConnection) {
	f.connection = connection
}

func (r *RedisGateway) SetConnection(pool RedisConnection) {
	r.pool = pool
}
