package broadcaster

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/rs/zerolog/log"

	"go-feedmaker/entity"
)

type (
	generationOut struct {
		ID            string  `json:"id"`
		Type          string  `json:"type"`
		Progress      uint    `json:"progress"`
		DataFetched   bool    `json:"data_fetched"`
		FilesUploaded uint    `json:"files_uploaded"`
		IsCanceled    bool    `json:"is_canceled"`
		StartTime     string  `json:"start_time"`
		EndTime       *string `json:"end_time"`
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

func marshalGeneration(generation *entity.Generation, w io.Writer) error {
	generationOut := makeGenerationOut(generation)
	encoder := json.NewEncoder(w)
	return encoder.Encode(generationOut)
}

func makeGenerationOut(generation *entity.Generation) *generationOut {
	generationOut := &generationOut{
		ID:            generation.ID,
		Type:          generation.Type,
		Progress:      generation.Progress,
		DataFetched:   generation.DataFetched,
		FilesUploaded: generation.FilesUploaded,
		IsCanceled:    generation.IsCanceled,
		StartTime:     formatTime(generation.StartTime),
	}
	if !generation.EndTime.IsZero() {
		endTime := formatTime(generation.EndTime)
		generationOut.EndTime = &endTime
	}
	return generationOut
}

func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
