package interactor

import (
	"context"
	"io"
	"sync"
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

	DataFetcher interface {
		StreamData(ctx context.Context) error
		OnDataFetched(func())
		OnProgress(func(progress uint))
	}

	FeedFactory interface {
		CreateDataFetcher(outStream chan<- []string) DataFetcher
		CreateFileFormatter(inStream <-chan []string, outStream chan<- io.ReadCloser) FileFormatter
		CreateUploader(inStream <-chan io.ReadCloser) Uploader
	}

	FileFormatter interface {
		FormatFiles(ctx context.Context) error
	}

	Uploader interface {
		UploadFiles(ctx context.Context) error
		OnUpload(func(uploadedNum uint))
	}

	FeedRepo interface {
		GetFactoryByGenerationType(generationType string) (FeedFactory, error)
		StoreGeneration(ctx context.Context, generation *entity.Generation) (*entity.Generation, error)
		UpdateGenerationState(ctx context.Context, generation *entity.Generation) error
		ListGenerations(ctx context.Context) ([]*entity.Generation, error)
		ListAllowedTypes(ctx context.Context) ([]string, error)
		IsAllowedType(ctx context.Context, generationType string) (bool, error)
		IsCanceled(ctx context.Context, generationID string) (bool, error)
		CancelGeneration(ctx context.Context, id string) error
		OnGenerationCanceled(ctx context.Context, id string, callback func()) error
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
	defer cancelCtx()
	go i.onGenerationCanceled(ctx, generation.ID, cancelCtx)
	errStream := make(chan error)
	recordStream := make(chan []string)
	fileStream := make(chan io.ReadCloser)

	dataFetcher := factory.CreateDataFetcher(recordStream)
	fileFormatter := factory.CreateFileFormatter(recordStream, fileStream)
	uploader := factory.CreateUploader(fileStream)
	dataFetcher.OnDataFetched(i.onDataFetched(generation))
	dataFetcher.OnProgress(i.onProgress(generation))
	uploader.OnUpload(i.onFileUploaded(generation))

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		defer close(recordStream)
		if err := dataFetcher.StreamData(ctx); err != nil {
			errStream <- err
		}
	}()
	go func() {
		defer wg.Done()
		defer close(fileStream)
		if err := fileFormatter.FormatFiles(ctx); err != nil {
			errStream <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := uploader.UploadFiles(ctx); err != nil {
			errStream <- err
		}
	}()
	go func() {
		wg.Wait()
		close(errStream)
	}()

	for err := range errStream {
		cancelCtx()
		return i.presenter.PresentErr(err)
	}
	return nil
}

func (i *feedInteractor) onProgress(generation *entity.Generation) func(uint) {
	return func(progress uint) {
		generation.SetProgress(progress)
		if err := i.feeds.UpdateGenerationState(context.Background(), generation); err != nil {
			log.Error().Err(err).
				Msgf("Cannot update progress for %s %v", generation.ID, progress)
		}
	}
}

func (i *feedInteractor) onFileUploaded(generation *entity.Generation) func(uint) {
	return func(uploadedNum uint) {
		generation.FilesUploaded++
		if err := i.feeds.UpdateGenerationState(context.Background(), generation); err != nil {
			log.Error().Err(err).
				Msgf("Cannot update file uploaded for %s %v", generation.ID, generation.FilesUploaded)
		}
	}
}

func (i *feedInteractor) onDataFetched(generation *entity.Generation) func() {
	return func() {
		generation.DataFetched = true
		if err := i.feeds.UpdateGenerationState(context.Background(), generation); err != nil {
			log.Error().Err(err).
				Msgf("Cannot update data fetched for %s", generation.ID)
		}
	}
}

func (i *feedInteractor) onGenerationCanceled(ctx context.Context, generationID string, callback func()) {
	err := i.feeds.OnGenerationCanceled(ctx, generationID, callback)
	if err != nil {
		log.Error().Err(err).
			Msgf("Cannot check if generation with id %s canceled", generationID)
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
