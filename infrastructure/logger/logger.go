package logger

import (
	"io"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Level      string `config:"level"`
	JsonOutput bool   `config:"json_output"`
}

func (c *Config) Apply() error {
	level, err := c.parseLevel()
	if err != nil {
		return err
	}
	c.applyLevel(level)

	if !c.JsonOutput {
		textWriter := zerolog.NewConsoleWriter()
		c.applyWriter(textWriter)
	}

	return nil
}

func (c *Config) parseLevel() (zerolog.Level, error) {
	level, err := zerolog.ParseLevel(strings.ToLower(c.Level))
	if err != nil {
		return zerolog.Disabled, err
	}
	return level, nil
}

func (c *Config) applyLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
}

func (c *Config) applyWriter(w io.Writer) {
	log.Logger = log.Output(w)
}
