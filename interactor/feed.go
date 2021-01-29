package interactor

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"go-feedmaker/entity"
)

type (
	FeedInteractor interface {
		Generate(ctx context.Context, generationType string) error
		ListGenerations(ctx context.Context) (interface{}, error)
		ListGenerationTypes(ctx context.Context) (interface{}, error)
		CancelGeneration(ctx context.Context, id string) error
	}

	Uploader interface {
		Upload(ctx context.Context, filepath string, file io.Reader) error
	}

	DataFetcher interface {
		FetchDataStream(ctx context.Context) (<-chan []string, error)
	}

	FeedFactory interface {
		CreateDataFetcher() DataFetcher
		CreateFileFormatter(dataStream <-chan []string) FileFormatter
	}

	FileFormatter interface {
		FormatFiles(ctx context.Context) (<-chan io.Reader, error)
	}

	FeedRepo interface {
		GetFactoryByGenerationType(generationType string) (FeedFactory, error)
		StoreGeneration(ctx context.Context, generation *entity.Generation) (*entity.Generation, error)
		ListGenerations(ctx context.Context) ([]*entity.Generation, error)
		ListAllowedTypes(ctx context.Context) ([]string, error)
		IsAllowedType(ctx context.Context, generationType string) (bool, error)
		CancelGeneration(ctx context.Context, id string) error
	}

	Presenter interface {
		PresentGenerationTypes([]string) interface{}
		PresentListGenerations(out *ListGenerationsOut) interface{}
		PresentErr(err error) error
	}

	feedInteractor struct {
		uploader  Uploader
		feeds     FeedRepo
		presenter Presenter
	}

	GenerationsOut struct {
		ID        string
		Type      string
		Progress  uint
		StartTime time.Time
		EndTime   time.Time
	}

	ListGenerationsOut []*GenerationsOut
)

func NewFeedInteractor(
	files Uploader,
	generations FeedRepo,
	presenter Presenter,
) *feedInteractor {
	return &feedInteractor{
		uploader:  files,
		feeds:     generations,
		presenter: presenter,
	}
}

func (i *feedInteractor) Generate(ctx context.Context, generationType string) error {
	factory, err := i.feeds.GetFactoryByGenerationType(generationType)
	if err != nil {
		return i.presenter.PresentErr(err)
	}

	generation, err := i.feeds.StoreGeneration(ctx, &entity.Generation{
		ID:        uuid.New().String(),
		Type:      generationType,
		StartTime: time.Now(),
	})
	if err != nil {
		return i.presenter.PresentErr(err)
	}

	dataFetcher := factory.CreateDataFetcher()
	dataStream, err := dataFetcher.FetchDataStream(ctx)
	if err != nil {
		return i.presenter.PresentErr(err)
	}
	fileFormatter := factory.CreateFileFormatter(dataStream)
	fileStream, err := fileFormatter.FormatFiles(ctx)
	if err != nil {
		return i.presenter.PresentErr(err)
	}

	for file := range fileStream {
		if err := i.uploader.Upload(ctx, generationType, file); err != nil {
			log.Error().Msgf("Cannot upload file, %s", err.Error())
		}
	}
	_ = generation
	return nil
}

func (i *feedInteractor) ListGenerations(ctx context.Context) (interface{}, error) {
	generations, err := i.feeds.ListGenerations(ctx)
	if err != nil {
		return nil, i.presenter.PresentErr(err)
	}
	return i.presenter.PresentListGenerations(makeListGenerationsOut(generations)), nil
}

func makeListGenerationsOut(generations []*entity.Generation) *ListGenerationsOut {
	out := ListGenerationsOut{}
	for _, generation := range generations {
		out = append(out, &GenerationsOut{
			ID:        generation.ID,
			Type:      generation.Type,
			Progress:  generation.Progress,
			StartTime: generation.StartTime,
			EndTime:   generation.EndTime,
		})
	}
	return &out
}

func (i *feedInteractor) ListGenerationTypes(ctx context.Context) (interface{}, error) {
	allowedTypes, err := i.feeds.ListAllowedTypes(ctx)
	if err != nil {
		return nil, i.presenter.PresentErr(err)
	}
	return i.presenter.PresentGenerationTypes(allowedTypes), nil
}

func (i *feedInteractor) CancelGeneration(ctx context.Context, id string) error {
	if err := i.feeds.CancelGeneration(ctx, id); err != nil {
		return i.presenter.PresentErr(err)
	}
	return nil
}
