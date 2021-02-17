package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/inhies/go-bytesize"
	"github.com/jlaffaye/ftp"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"

	"go-feedmaker/adapter/presenter"
	"go-feedmaker/adapter/repository"
	"go-feedmaker/entity"
	"go-feedmaker/infrastructure/config"
	"go-feedmaker/infrastructure/gateway"
	"go-feedmaker/infrastructure/rest"
	"go-feedmaker/infrastructure/rest/broadcaster"
	"go-feedmaker/infrastructure/scheduler"
	"go-feedmaker/interactor"
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
		log.Fatal().Err(err).Msg("Can't configure logger")
	}

	redisGateway := &gateway.RedisGateway{
		Config: conf.Redis,
		Dialer: new(redisDialer),
	}
	if err := redisGateway.Connect(); err != nil {
		log.Fatal().Err(err).Msg("Can't connect to Redis")
	}
	defer redisGateway.Disconnect()

	ftpGateway := &gateway.FtpGateway{
		Dialer: new(ftpDialer),
		Config: conf.Ftp,
	}
	if err := ftpGateway.Connect(); err != nil {
		log.Fatal().Err(err).Msg("Can't connect to Ftp")
	}
	defer ftpGateway.Disconnect()

	feedRepoConfig, err := initFeedRepoConfig(conf.Feeds)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't initialize config for feed repo")
	}
	sqlGateways := make([]*gateway.SqlGateway, 0, len(conf.Feeds))
	for key, conf := range conf.Feeds {
		sqlGateway := &gateway.SqlGateway{
			DriverName: conf.Database.Driver,
			DSN:        conf.Database.Dsn,
		}
		sqlGateways = append(sqlGateways, sqlGateway)
		if err := sqlGateway.Connect(); err != nil {
			log.Fatal().Err(err).Msgf("Can't connect to database for feed %s", key)
		}
		feedRepoConfig[key].SqlGateway = sqlGateway.DB()
	}
	defer closeSqlGateways(sqlGateways)

	feedRepo := repository.NewFeedRepo(feedRepoConfig, redisGateway, ftpGateway)
	feedPresenter := new(presenter.Presenter)
	feedInteractor := interactor.NewFeedInteractor(feedRepo, feedPresenter)

	scheduleSaver := scheduler.NewScheduleSaver(redisGateway)
	taskScheduler := scheduler.New(cron.New(), scheduleSaver)
	taskScheduler.Start()
	defer taskScheduler.Stop()
	if err := taskScheduler.ScheduleAllSavedGenerations(feedInteractor); err != nil {
		log.Fatal().Err(err).Msg("can't schedule saved generations")
	}
	handler := rest.NewHandler(feedInteractor, taskScheduler)

	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	progressBroadcaster := broadcaster.NewBroadcaster()
	go progressBroadcaster.Start()
	defer progressBroadcaster.Stop()
	wsHandler := rest.NewWSHandler(upgrader, progressBroadcaster)

	generationsProgress := make(chan *entity.Generation)
	go feedInteractor.WatchGenerationsProgress(context.Background(), generationsProgress)
	go progressBroadcaster.BroadcastGenerationsProgress(generationsProgress)

	router := rest.NewRouter(handler, wsHandler)
	apiServer := rest.NewAPIServer(&conf.Api, router)
	go apiServer.Start()
	defer apiServer.Stop()

	log.Info().Msgf("Successfully started server")
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info().Msgf("Server was stopped")
}

func initFeedRepoConfig(config map[string]config.FeedConfig) (map[string]*repository.FeedConfig, error) {
	res := make(map[string]*repository.FeedConfig, len(config))
	for key, conf := range config {
		countQuery, err := readSqlFromFile(conf.CountQueryFilename)
		if err != nil {
			return nil, err
		}
		selectQuery, err := readSqlFromFile(conf.SelectQueryFilename)
		if err != nil {
			return nil, err
		}
		fileSizeLimit, err := bytesize.Parse(conf.FileSizeLimit)
		if err != nil {
			return nil, err
		}

		res[key] = &repository.FeedConfig{
			CountQuery:    countQuery,
			SelectQuery:   selectQuery,
			FileSizeLimit: fileSizeLimit,
			FileLineLimit: conf.FileLineLimit,
		}
	}
	return res, nil
}

func closeSqlGateways(sqlGateways []*gateway.SqlGateway) {
	for _, sqlGateway := range sqlGateways {
		if err := sqlGateway.Disconnect(); err != nil {
			log.Error().Err(err).Msg("Can't close sql gateway")
		}
	}
}

func readSqlFromFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	sql, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(sql), nil
}
