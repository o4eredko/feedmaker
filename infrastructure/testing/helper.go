package helper

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func OpenFile(t *testing.T, filename string) *os.File {
	file, err := os.Open(filename)
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, file.Close())
	})
	return file
}

func TimeoutCtx(t *testing.T, parent context.Context, timeout time.Duration) context.Context {
	ctx, cancel := context.WithTimeout(parent, timeout)
	t.Cleanup(cancel)
	return ctx
}
