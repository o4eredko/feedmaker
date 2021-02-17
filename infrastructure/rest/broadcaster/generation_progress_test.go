package broadcaster_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-feedmaker/entity"
	"go-feedmaker/infrastructure/rest/broadcaster"
)

func Test_broadcaster_BroadcastGenerationsProgress(t *testing.T) {
	fields := defaultBroadcasterFields()
	b := broadcaster.NewBroadcaster()
	b.SetBroadcast(fields.broadcast)
	generationProgress := make(chan *entity.Generation)
	go b.BroadcastGenerationsProgress(generationProgress)
	generationProgress <- defaultGeneration
	gotMessage := <-fields.broadcast
	buf := new(bytes.Buffer)
	assert.NoError(t, broadcaster.MarshalGeneration(defaultGeneration, buf))
	assert.Equal(t, buf.Bytes(), gotMessage)
}
