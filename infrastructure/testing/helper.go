package helper

import (
	"os"
	"testing"

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
