package main

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"

	"go-feedmaker/adapter/repository"
	"go-feedmaker/infrastructure/gateway"
)

type redisDialer struct{}

func (r *redisDialer) Dial(network, addr string, options ...redis.DialOption) (gateway.RedisConnection, error) {
	return redis.Dial(network, addr, options...)
}

// func main() {
// 	redisGateway := &gateway.RedisGateway{
// 		Config: gateway.RedisConfig{
// 			Host:        "localhost",
// 			Port:        "6379",
// 			ConnTimeout: time.Second,
// 		},
// 		Dialer: new(redisDialer),
// 	}
// 	if err := redisGateway.Connect(); err != nil {
// 		panic(err)
// 	}
// 	defer redisGateway.Disconnect()
//
// 	repo := repository.NewFeedRepo(redisGateway)
//
// 	generation := &entity.Generation{
// 		ID:        "abc",
// 		Type:      "test type",
// 		Progress:  50,
// 		StartTime: time.Now(),
// 	}
//
// 	// go func() {
// 	// 	time.Sleep(time.Second * 3)
// 	// 	if err := repo.CancelGeneration(context.Background(), generation.ID); err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// }()
// 	callback := func() {
// 		log.Info().Msgf("Generation %s was canceled", generation.ID)
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	if err := repo.OnGenerationCanceled(ctx, generation.ID, callback); err != nil {
// 		panic(err)
// 	}
// 	log.Info().Msgf("returned")
// }
var (
	username = "feed_maker"
	password = "ohRS-Tx6d_O56j3AzMRW"
	host     = "warehouse.jo"
	database = "Marketing"
	dsn      = fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", username, password, host, database)
)

func main() {
	// recordStream := make(chan []string)
	sqlGateway := gateway.SqlGateway{DriverName: "mssql", DSN: dsn}
	if err := sqlGateway.Connect(); err != nil {
		panic(err)
	}
	defer sqlGateway.Disconnect()

	countQuery := "SELECT COUNT(*) FROM dbo.accounts"
	selectQuery := "SELECT id, original_id, name, created_at FROM dbo.accounts"
	fetcher := repository.SqlDataFetcher{
		// OutStream:   recordStream,
		SelectQuery: selectQuery,
		CountQuery:  countQuery,
		Db:          sqlGateway.DB(),
	}
	count, err := fetcher.CountRecords(context.Background())
	if err != nil {
		panic(err)
	}
	log.Info().Msgf("%d", count)
	// go func() {
	// 	if err := fetcher.StreamData(context.Background()); err != nil {
	// 		panic(err)
	// 	}
	// 	close(recordStream)
	// }()
	// for record := range recordStream {
	// 	log.Info().Msgf("%s", record)
	// }
}
