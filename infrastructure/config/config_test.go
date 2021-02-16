package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
				assert.NotNil(t, c)
			}
			if tc.wantPanic {
				assert.Panics(t, testFunc)
			} else {
				assert.NotPanics(t, testFunc)
			}
		})
	}
}
