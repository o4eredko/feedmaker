package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/types/objectpath"

	"go-feedmaker/infrastructure/config"
)

func TestLoadFromFile(t *testing.T) {
	testCases := []struct {
		name      string
		filename  objectpath.Path
		wantPanic bool
	}{
		{
			name:      "succeed",
			filename:  config.DefaultConfig,
			wantPanic: false,
		},
		{
			name:      "no such file",
			filename:  "invalid.json",
			wantPanic: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testFunc := func() {
				c := config.LoadFromFile(tc.filename)
				require.NotNil(t, c)
				assert.NotEmpty(t, c.Logger)
				assert.NotEmpty(t, c.Ftp)
				assert.NotEmpty(t, c.Redis)
				assert.NotEmpty(t, c.Feeds)
			}
			if tc.wantPanic {
				assert.Panics(t, testFunc)
			} else {
				assert.NotPanics(t, testFunc)
			}
		})
	}
}
