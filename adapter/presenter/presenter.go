package presenter

import (
	"fmt"
	"time"

	"go-feedmaker/interactor"
)

type (
	Presenter struct{}

	generationOut struct {
		ID            string  `json:"id"`
		Type          string  `json:"type"`
		Progress      string  `json:"progress"`
		DataFetched   bool    `json:"data_fetched"`
		FilesUploaded uint    `json:"files_uploaded"`
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
		Progress:      fmt.Sprintf("%d%%", generation.Progress),
		DataFetched:   generation.DataFetched,
		FilesUploaded: generation.FilesUploaded,
		StartTime:     generation.StartTime.UTC().Format(time.RFC3339),
	}
	if !generation.EndTime.IsZero() {
		endTime := generation.EndTime.UTC().Format(time.RFC3339)
		generationOut.EndTime = &endTime
	}
	return generationOut
}
