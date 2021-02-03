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
		GenerateFeed(ctx context.Context, generationType string) error
		ListGenerations(ctx context.Context) (interface{}, error)
		ListGenerationTypes(ctx context.Context) (interface{}, error)
		CancelGeneration(ctx context.Context, id string) error
	}

	Uploader interface {
		Upload(ctx context.Context, filepath string, file io.Reader) error
	}

	DataFetcher interface {
		CountRecords(ctx context.Context) (uint, error)
		StreamData(ctx context.Context) error
	}

	FeedFactory interface {
		CreateDataFetcher(outStream chan<- []string) DataFetcher
		CreateFileFormatter(inStream <-chan []string, outStream chan<- io.ReadCloser) FileFormatter
	}

	FileFormatter interface {
		FormatFiles(ctx context.Context) error
	}

	FeedRepo interface {
		GetFactoryByGenerationType(generationType string) (FeedFactory, error)
		StoreGeneration(ctx context.Context, generation *entity.Generation) (*entity.Generation, error)
		UpdateProgress(ctx context.Context, generationID string, progress int) error
		ListGenerations(ctx context.Context) ([]*entity.Generation, error)
		ListAllowedTypes(ctx context.Context) ([]string, error)
		IsAllowedType(ctx context.Context, generationType string) (bool, error)
		IsCanceled(ctx context.Context, generationID string) (bool, error)
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

	GenerationsOut entity.Generation

	ListGenerationsOut []*GenerationsOut
)

func NewFeedInteractor(files Uploader, generations FeedRepo, presenter Presenter) *feedInteractor {
	return &feedInteractor{
		uploader:  files,
		feeds:     generations,
		presenter: presenter,
	}
}

func (i *feedInteractor) GenerateFeed(ctx context.Context, generationType string) error {
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

	ctx, cancelCtx := context.WithCancel(ctx)
	go i.onGenerationCanceled(ctx, generation.ID, cancelCtx)
	errStream := make(chan error)
	recordStream := make(chan []string)
	fileStream := make(chan io.ReadCloser)
	defer cancelCtx()
	defer close(errStream)

	dataFetcher := factory.CreateDataFetcher(recordStream)
	fileFormatter := factory.CreateFileFormatter(recordStream, fileStream)

	go i.fetchData(ctx, dataFetcher, errStream)
	go i.formatFiles(ctx, fileFormatter, errStream)
	go i.uploadFiles(ctx, fileStream, generationType, errStream)
	for err := range errStream {
		cancelCtx()
		return i.presenter.PresentErr(err)
	}
	return nil
}

func (i *feedInteractor) uploadFiles(ctx context.Context, fileStream <-chan io.ReadCloser, generationType string, errStream chan<- error) {
	for file := range fileStream {
		if err := i.uploader.Upload(ctx, generationType, file); err != nil {
			errStream <- err
		}
		if err := file.Close(); err != nil {
			log.Error().Err(err).Msgf("Cannot close file during %s generation", generationType)
		}
	}
}

func (i *feedInteractor) fetchData(ctx context.Context, fetcher DataFetcher, errStream chan<- error) {
	if err := fetcher.StreamData(ctx); err != nil {
		errStream <- err
	}
}

func (i *feedInteractor) formatFiles(ctx context.Context, formatter FileFormatter, errStream chan<- error) {
	if err := formatter.FormatFiles(ctx); err != nil {
		errStream <- err
	}
}

func (i *feedInteractor) onProgress(generationID string, progress int) {
	err := i.feeds.UpdateProgress(context.Background(), generationID, progress)
	if err != nil {
		log.Error().Err(err).Msgf(
			"Cannot set progress for generation = %s progress = %d",
			generationID, progress,
		)
	}
}

func (i *feedInteractor) onGenerationCanceled(ctx context.Context, generationID string, f func()) {
	for {
		isRejected, err := i.feeds.IsCanceled(ctx, generationID)
		if err != nil {
			log.Error().
				Err(err).
				Msgf("Cannot check if generation with id %s rejected", generationID)
		} else if isRejected {
			f()
		}
		time.Sleep(time.Second)
	}
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
