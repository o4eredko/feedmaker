package interactor

import (
	"context"
	"io"

	"go-feedmaker/entity"
)

type (
	FeedInteractor interface {
		Generate(ctx context.Context, generationType string) (interface{}, error)
		ListGenerations(ctx context.Context) (interface{}, error)
		ListGenerationTypes(ctx context.Context) (interface{}, error)
		PauseGeneration(ctx context.Context) error
		ResumeGeneration(ctx context.Context) error
	}

	FileRepo interface {
		Upload(ctx context.Context, file io.Reader) error
	}

	CsvRepo interface {
		io.Reader
	}

	GenerationRepo interface {
		Store(ctx context.Context, generation *entity.Generation) error
		List(ctx context.Context) ([]*entity.Generation, error)
		ListAllowedTypes(ctx context.Context) ([]string, error)
		IsAllowedType(ctx context.Context, generationType string) (bool, error)
	}

	Presenter interface {
		PresentGenerationTypes([]string) interface{}
		PresentErr(err error) error
	}

	feedInteractor struct {
		files       FileRepo
		generations GenerationRepo
		csvFetcher  CsvRepo
		presenter   Presenter
	}
)

func NewFeedInteractor(
	csvFetcher CsvRepo,
	files FileRepo,
	generations GenerationRepo,
	presenter Presenter,
) *feedInteractor {
	return &feedInteractor{
		csvFetcher:  csvFetcher,
		files:       files,
		generations: generations,
		presenter:   presenter,
	}
}

func (f *feedInteractor) Generate(ctx context.Context, generationType string) (interface{}, error) {
	isAllowed, err := f.generations.IsAllowedType(ctx, generationType)
	if err != nil {
		return nil, f.presenter.PresentErr(err)
	} else if !isAllowed {
		return nil, f.presenter.PresentErr(entity.ErrInvalidGenerationType)
	}
	return nil, nil
}

func (f *feedInteractor) ListGenerations(ctx context.Context) (interface{}, error) {
	panic("implement me")
}

func (f *feedInteractor) ListGenerationTypes(ctx context.Context) (interface{}, error) {
	allowedTypes, err := f.generations.ListAllowedTypes(ctx)
	if err != nil {
		return nil, f.presenter.PresentErr(err)
	}
	return f.presenter.PresentGenerationTypes(allowedTypes), nil
}

func (f *feedInteractor) PauseGeneration(ctx context.Context) error {
	panic("implement me")
}

func (f *feedInteractor) ResumeGeneration(ctx context.Context) error {
	panic("implement me")
}
