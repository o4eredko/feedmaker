package broadcaster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"go-feedmaker/entity"
)

type (
	generationOut struct {
		ID            string `json:"id"`
		Type          string `json:"type"`
		Progress      string `json:"progress"`
		DataFetched   bool   `json:"data_fetched"`
		FilesUploaded uint   `json:"files_uploaded"`
		StartTime     string `json:"start_time"`
		EndTime       string `json:"end_time"`
	}
)

func (b *broadcaster) BroadcastGenerationsProgress(generationsProgress <-chan *entity.Generation) {
	for generation := range generationsProgress {
		b.broadcastGeneration(generation)
	}
}

func (b *broadcaster) broadcastGeneration(generation *entity.Generation) {
	buf := new(bytes.Buffer)
	if err := marshalGeneration(generation, buf); err != nil {
		log.Error().
			Err(err).
			Interface("generation", generation).
			Msg("generation marshal")
		return
	}
	b.Broadcast(buf.Bytes())
}

func marshalGeneration(generation *entity.Generation, buf *bytes.Buffer) error {
	generationOut := makeGenerationOut(generation)
	encoder := json.NewEncoder(buf)
	return encoder.Encode(generationOut)
}

func makeGenerationOut(generation *entity.Generation) *generationOut {
	generationOut := &generationOut{
		ID:            generation.ID,
		Type:          generation.Type,
		Progress:      fmt.Sprintf("%d%%", generation.Progress),
		DataFetched:   generation.DataFetched,
		FilesUploaded: generation.FilesUploaded,
		StartTime:     generation.StartTime.UTC().Format(time.RFC3339),
		EndTime:       generation.EndTime.UTC().Format(time.RFC3339),
	}
	return generationOut
}