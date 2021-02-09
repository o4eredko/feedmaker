package gateway

func (f *FtpGateway) SetConnection(connection FtpConnection) {
	f.connection = connection
}

func (r *RedisGateway) SetConnection(connection RedisConnection) {
	r.connection = connection
	r.RedisClient = connection
}
