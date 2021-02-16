package config

import (
	"path"
	"runtime"

	"github.com/o4eredko/configuro"
	"golang.org/x/tools/go/types/objectpath"

	"go-feedmaker/infrastructure/gateway"
	"go-feedmaker/infrastructure/logger"
)

type Config struct {
	Logger logger.Config
	Redis  gateway.RedisConfig
	Ftp    gateway.FtpConfig
}

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
