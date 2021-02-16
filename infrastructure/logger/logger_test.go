package logger_test

import (
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"go-feedmaker/infrastructure/logger"
)

func TestConfig_Apply(t *testing.T) {
	tests := []struct {
		name    string
		config  logger.Config
		wantErr bool
	}{
		{
			name:   "succeed",
			config: logger.Config{Level: "debug"},
		},
		{
			name:   "succeed",
			config: logger.Config{Level: "Error", JsonOutput: true},
		},
		{
			name:   "succeed",
			config: logger.Config{Level: "INFO"},
		},
		{
			name:   "succeed",
			config: logger.Config{Level: "wArN"},
		},
		{
			name:   "succeed",
			config: logger.Config{Level: "Error"},
		},
		{
			name:    "invalid level",
			config:  logger.Config{Level: "foobar"},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotErr := tc.config.Apply()
			if tc.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)

				wantLevel := strings.ToLower(tc.config.Level)
				gotLevel := zerolog.GlobalLevel().String()
				assert.Equal(t, wantLevel, gotLevel)
			}
		})
	}
}
