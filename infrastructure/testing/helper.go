package helper

import (
	"context"
	"encoding/csv"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func OpenFile(t *testing.T, filename string) *os.File {
	file, err := os.Open(filename)
	require.NoError(t, err)
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

func ReadCsvFromFile(t *testing.T, filename string) [][]string {
	records, err := csv.NewReader(OpenFile(t, filename)).ReadAll()
	require.NoError(t, err)
	return records
}
