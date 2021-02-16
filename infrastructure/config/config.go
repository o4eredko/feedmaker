package config

import (
	"path"
	"runtime"

	"github.com/o4eredko/configuro"
	"golang.org/x/tools/go/types/objectpath"

	"go-feedmaker/infrastructure/gateway"
	"go-feedmaker/infrastructure/logger"
	"go-feedmaker/infrastructure/rest"
)

type (
	FeedConfig struct {
		CountQueryFilename  string `config:"count_query"`
		SelectQueryFilename string `config:"select_query"`
		FileSizeLimit       string `config:"size_limit"`
		FileLineLimit       uint   `config:"line_limit"`
		Database            struct {
			Driver string
			Dsn    string
		}
	}

	Config struct {
		Logger logger.Config
		Redis  gateway.RedisConfig
		Ftp    gateway.FtpConfig
		Feeds  map[string]FeedConfig
		Api    rest.Config
	}
)

func LoadFromFile(filename objectpath.Path) *Config {
	configPath := path.Join(getPackageDir(), string(filename))
	configLoader, err := configuro.NewConfig(
		configuro.WithLoadFromConfigFile(configPath, true),
	)
	if err != nil {
		panic(err)
	}
	config := new(Config)
	if err := configLoader.Load(config); err != nil {
		panic(err)
	}
	return config
}

func getPackageDir() string {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	return path.Dir(currentFilePath)
}
