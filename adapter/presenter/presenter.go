package presenter

import (
	"time"

	"go-feedmaker/interactor"
)

type (
	Presenter struct{}

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

func (p *Presenter) PresentGenerationTypes(out []string) interface{} {
	return out
}

func (p *Presenter) PresentListGenerations(generations *interactor.ListGenerationsOut) interface{} {
	generationsOut := make([]*generationOut, len(*generations))
	for i, generation := range *generations {
		generationOut := makeGenerationOut(generation)
		generationsOut[i] = generationOut
	}
	return generationsOut
}

func (p *Presenter) PresentErr(err error) error {
	return err
}

func makeGenerationOut(generation *interactor.GenerationsOut) *generationOut {
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
