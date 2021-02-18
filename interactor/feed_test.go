package interactor_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-feedmaker/entity"
	"go-feedmaker/interactor"
	"go-feedmaker/interactor/mocks"
)

var (
	defaultID      = "qwerwqerjwejr"
	defaultErr     = errors.New("default test error")
	errPassThrough = func(err error) error {
		return err
	}
	generation1 = entity.Generation{
		ID:        "hesoyam",
		Type:      "test",
		Progress:  43,
		StartTime: time.Now(),
	}
	generation2 = entity.Generation{
		ID:        "qwerty",
		Type:      "test",
		Progress:  100,
		StartTime: time.Now(),
		EndTime:   time.Now(),
	}
)

type fields struct {
	uploader      *mocks.Uploader
	feeds         *mocks.FeedRepo
	presenter     *mocks.Presenter
	factory       *mocks.FeedFactory
	fileFormatter *mocks.FileFormatter
	dataFetcher   *mocks.DataFetcher
}

func defaultFields() *fields {
	return &fields{
		uploader:      new(mocks.Uploader),
		feeds:         new(mocks.FeedRepo),
		presenter:     new(mocks.Presenter),
		factory:       new(mocks.FeedFactory),
		fileFormatter: new(mocks.FileFormatter),
		dataFetcher:   new(mocks.DataFetcher),
	}
}

func (f *fields) newInteractor() interactor.FeedInteractor {
	return interactor.NewFeedInteractor(f.feeds, f.presenter)
}

func (f *fields) assertExpectations(t *testing.T) {
	f.uploader.AssertExpectations(t)
	f.feeds.AssertExpectations(t)
	f.presenter.AssertExpectations(t)
	f.factory.AssertExpectations(t)
	f.fileFormatter.AssertExpectations(t)
	f.dataFetcher.AssertExpectations(t)
}

func TestNewFeedInteractor(t *testing.T) {
	fields := defaultFields()
	i := interactor.NewFeedInteractor(fields.feeds, fields.presenter)
	assert.Equal(t, fields.feeds, i.GenerationRepo())
	assert.Equal(t, fields.presenter, i.Presenter())
}

func TestFeedInteractor_GenerateFeed(t *testing.T) {
	type args struct {
		ctx            context.Context
		generationType string
	}
	defaultArgs := func() *args {
		return &args{
			ctx:            context.Background(),
			generationType: "test",
		}
	}
	testCases := []struct {
		name       string
		args       *args
		setupMocks func(*args, *fields)
		wantErr    error
	}{
		{
			name: "succeed",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				generationMatches := func(g *entity.Generation) bool {
					timeIsAlmostEqual := g.StartTime.Sub(time.Now()) < time.Second
					return g.Progress == 0 && timeIsAlmostEqual && g.Type == a.generationType && len(g.ID) > 0
				}

				f.feeds.On("GetFactoryByGenerationType", a.generationType).
					Return(f.factory, nil)
				f.feeds.On("StoreGeneration", a.ctx, mock.MatchedBy(generationMatches)).
					Return(nil)
				f.feeds.On("OnGenerationCanceled", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				f.factory.On("CreateDataFetcher", mock.Anything).Return(f.dataFetcher)
				f.factory.On("CreateFileFormatter", mock.Anything, mock.Anything).Return(f.fileFormatter)
				f.factory.On("CreateUploader", mock.Anything).Return(f.uploader)

				f.dataFetcher.
					On("StreamData", mock.Anything).Return(nil).
					On("OnDataFetched", mock.Anything).Return(nil).
					On("OnProgress", mock.Anything).Return(nil)
				f.fileFormatter.
					On("FormatFiles", mock.Anything).Return(nil)
				f.uploader.
					On("UploadFiles", mock.Anything).Return(nil).
					On("OnUpload", mock.Anything).Return(nil)
			},
		},
		{
			name: "unknown generation type",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				f.feeds.On("GetFactoryByGenerationType", a.generationType).Return(nil, defaultErr)
				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)
			},
			wantErr: defaultErr,
		},
		{
			name: "store generation error",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				generationMatches := func(g *entity.Generation) bool {
					timeIsAlmostEqual := g.StartTime.Sub(time.Now()) < time.Second
					return g.Progress == 0 && timeIsAlmostEqual && g.Type == a.generationType && len(g.ID) > 0
				}
				f.feeds.On("GetFactoryByGenerationType", a.generationType).Return(f.factory, nil)
				f.feeds.On("StoreGeneration", a.ctx, mock.MatchedBy(generationMatches)).Return(defaultErr)
				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)
			},
			wantErr: defaultErr,
		},
		{
			name: "stream data error",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				generationMatches := func(g *entity.Generation) bool {
					timeIsAlmostEqual := g.StartTime.Sub(time.Now()) < time.Second
					return g.Progress == 0 && timeIsAlmostEqual && g.Type == a.generationType && len(g.ID) > 0
				}

				f.feeds.On("GetFactoryByGenerationType", a.generationType).
					Return(f.factory, nil)
				f.feeds.On("StoreGeneration", a.ctx, mock.MatchedBy(generationMatches)).
					Return(nil)
				f.feeds.On("OnGenerationCanceled", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				f.factory.On("CreateDataFetcher", mock.Anything).Return(f.dataFetcher)
				f.factory.On("CreateFileFormatter", mock.Anything, mock.Anything).Return(f.fileFormatter)
				f.factory.On("CreateUploader", mock.Anything).Return(f.uploader)

				f.dataFetcher.
					On("StreamData", mock.Anything).Return(defaultErr).After(time.Millisecond*5).
					On("OnDataFetched", mock.Anything).Return(nil).
					On("OnProgress", mock.Anything).Return(nil)
				f.fileFormatter.
					On("FormatFiles", mock.Anything).Return(nil)
				f.uploader.
					On("UploadFiles", mock.Anything).Return(nil).
					On("OnUpload", mock.Anything).Return(nil)

				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)
			},
			wantErr: defaultErr,
		},
		{
			name: "format files error",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				generationMatches := func(g *entity.Generation) bool {
					timeIsAlmostEqual := g.StartTime.Sub(time.Now()) < time.Second
					return g.Progress == 0 && timeIsAlmostEqual && g.Type == a.generationType && len(g.ID) > 0
				}

				f.feeds.On("GetFactoryByGenerationType", a.generationType).
					Return(f.factory, nil)
				f.feeds.On("StoreGeneration", a.ctx, mock.MatchedBy(generationMatches)).
					Return(nil)
				f.feeds.On("OnGenerationCanceled", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				f.factory.On("CreateDataFetcher", mock.Anything).Return(f.dataFetcher)
				f.factory.On("CreateFileFormatter", mock.Anything, mock.Anything).Return(f.fileFormatter)
				f.factory.On("CreateUploader", mock.Anything).Return(f.uploader)

				f.dataFetcher.
					On("StreamData", mock.Anything).Return(nil).
					On("OnDataFetched", mock.Anything).Return(nil).
					On("OnProgress", mock.Anything).Return(nil)
				f.fileFormatter.
					On("FormatFiles", mock.Anything).Return(defaultErr).After(time.Millisecond * 5)
				f.uploader.
					On("UploadFiles", mock.Anything).Return(nil).
					On("OnUpload", mock.Anything).Return(nil)

				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)
			},
			wantErr: defaultErr,
		},
		{
			name: "upload files error",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				generationMatches := func(g *entity.Generation) bool {
					timeIsAlmostEqual := g.StartTime.Sub(time.Now()) < time.Second
					return g.Progress == 0 && timeIsAlmostEqual && g.Type == a.generationType && len(g.ID) > 0
				}

				f.feeds.On("GetFactoryByGenerationType", a.generationType).
					Return(f.factory, nil)
				f.feeds.On("StoreGeneration", a.ctx, mock.MatchedBy(generationMatches)).
					Return(nil)
				f.feeds.On("OnGenerationCanceled", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				f.factory.On("CreateDataFetcher", mock.Anything).Return(f.dataFetcher)
				f.factory.On("CreateFileFormatter", mock.Anything, mock.Anything).Return(f.fileFormatter)
				f.factory.On("CreateUploader", mock.Anything).Return(f.uploader)

				f.dataFetcher.
					On("StreamData", mock.Anything).Return(nil).
					On("OnDataFetched", mock.Anything).Return(nil).
					On("OnProgress", mock.Anything).Return(nil)
				f.fileFormatter.
					On("FormatFiles", mock.Anything).Return(nil)
				f.uploader.
					On("UploadFiles", mock.Anything).Return(defaultErr).After(time.Millisecond*5).
					On("OnUpload", mock.Anything).Return(nil)

				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)
			},
			wantErr: defaultErr,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fields := defaultFields()
			interactor := fields.newInteractor()
			testCase.setupMocks(testCase.args, fields)

			gotErr := interactor.GenerateFeed(testCase.args.ctx, testCase.args.generationType)

			assert.Equal(t, testCase.wantErr, gotErr)
			fields.assertExpectations(t)
		})
	}
}

func TestFeedInteractor_RestartGeneration(t *testing.T) {
	type args struct {
		ctx          context.Context
		generationID string
	}
	defaultArgs := func() *args {
		return &args{
			ctx:          context.Background(),
			generationID: uuid.NewString(),
		}
	}
	testCases := []struct {
		name       string
		args       *args
		setupMocks func(*args, *fields)
		wantErr    error
	}{
		{
			name: "succeed",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				f.feeds.On("GetGeneration", a.ctx, a.generationID).
					Return(&entity.Generation{
						ID:            a.generationID,
						Type:          "test",
						Progress:      43,
						DataFetched:   true,
						FilesUploaded: 1,
						IsCanceled:    true,
						StartTime:     time.Unix(10, 0),
					}, nil)
				f.feeds.On("GetFactoryByGenerationType", "test").Return(f.factory, nil)
				f.feeds.On("UpdateGenerationState", a.ctx, &entity.Generation{
					ID:            a.generationID,
					Type:          "test",
					Progress:      0,
					DataFetched:   false,
					FilesUploaded: 0,
					IsCanceled:    false,
					StartTime:     time.Unix(10, 0),
				}).Return(nil)
				f.feeds.On("OnGenerationCanceled", mock.Anything, mock.Anything, mock.Anything).Return(nil)

				f.factory.On("CreateDataFetcher", mock.Anything).Return(f.dataFetcher)
				f.factory.On("CreateFileFormatter", mock.Anything, mock.Anything).Return(f.fileFormatter)
				f.factory.On("CreateUploader", mock.Anything).Return(f.uploader)

				f.dataFetcher.
					On("StreamData", mock.Anything).Return(nil).
					On("OnDataFetched", mock.Anything).Return(nil).
					On("OnProgress", mock.Anything).Return(nil)
				f.fileFormatter.
					On("FormatFiles", mock.Anything).Return(nil)
				f.uploader.
					On("UploadFiles", mock.Anything).Return(nil).
					On("OnUpload", mock.Anything).Return(nil)
			},
		},
		{
			name: "get generation error",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				f.feeds.On("GetGeneration", a.ctx, a.generationID).Return(nil, defaultErr)
				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)
			},
			wantErr: defaultErr,
		},
		{
			name: "unknown generation type",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				f.feeds.On("GetGeneration", a.ctx, a.generationID).
					Return(&entity.Generation{
						ID:            a.generationID,
						Type:          "test",
						Progress:      43,
						DataFetched:   true,
						FilesUploaded: 1,
						StartTime:     time.Unix(10, 0),
					}, nil)
				f.feeds.On("GetFactoryByGenerationType", "test").Return(nil, defaultErr)
				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)
			},
			wantErr: defaultErr,
		},
		{
			name: "stream data error",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				f.feeds.On("GetGeneration", a.ctx, a.generationID).
					Return(&entity.Generation{
						ID:            a.generationID,
						Type:          "test",
						Progress:      43,
						DataFetched:   true,
						FilesUploaded: 1,
						StartTime:     time.Unix(10, 0),
					}, nil)
				f.feeds.On("GetFactoryByGenerationType", "test").Return(f.factory, nil)
				f.feeds.On("UpdateGenerationState", a.ctx, &entity.Generation{
					ID:            a.generationID,
					Type:          "test",
					Progress:      0,
					DataFetched:   false,
					FilesUploaded: 0,
					StartTime:     time.Unix(10, 0),
				}).Return(nil)
				f.feeds.On("OnGenerationCanceled", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				f.factory.On("CreateDataFetcher", mock.Anything).Return(f.dataFetcher)
				f.factory.On("CreateFileFormatter", mock.Anything, mock.Anything).Return(f.fileFormatter)
				f.factory.On("CreateUploader", mock.Anything).Return(f.uploader)

				f.dataFetcher.
					On("StreamData", mock.Anything).Return(defaultErr).After(time.Millisecond*5).
					On("OnDataFetched", mock.Anything).Return(nil).
					On("OnProgress", mock.Anything).Return(nil)
				f.fileFormatter.
					On("FormatFiles", mock.Anything).Return(nil)
				f.uploader.
					On("UploadFiles", mock.Anything).Return(nil).
					On("OnUpload", mock.Anything).Return(nil)

				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)
			},
			wantErr: defaultErr,
		},
		{
			name: "format files error",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				f.feeds.On("GetGeneration", a.ctx, a.generationID).
					Return(&entity.Generation{
						ID:            a.generationID,
						Type:          "test",
						Progress:      43,
						DataFetched:   true,
						FilesUploaded: 1,
						StartTime:     time.Unix(10, 0),
					}, nil)
				f.feeds.On("GetFactoryByGenerationType", "test").Return(f.factory, nil)
				f.feeds.On("UpdateGenerationState", a.ctx, &entity.Generation{
					ID:            a.generationID,
					Type:          "test",
					Progress:      0,
					DataFetched:   false,
					FilesUploaded: 0,
					StartTime:     time.Unix(10, 0),
				}).Return(nil)
				f.feeds.On("OnGenerationCanceled", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				f.factory.On("CreateDataFetcher", mock.Anything).Return(f.dataFetcher)
				f.factory.On("CreateFileFormatter", mock.Anything, mock.Anything).Return(f.fileFormatter)
				f.factory.On("CreateUploader", mock.Anything).Return(f.uploader)

				f.dataFetcher.
					On("StreamData", mock.Anything).Return(nil).
					On("OnDataFetched", mock.Anything).Return(nil).
					On("OnProgress", mock.Anything).Return(nil)
				f.fileFormatter.
					On("FormatFiles", mock.Anything).Return(defaultErr).After(time.Millisecond * 5)
				f.uploader.
					On("UploadFiles", mock.Anything).Return(nil).
					On("OnUpload", mock.Anything).Return(nil)

				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)
			},
			wantErr: defaultErr,
		},
		{
			name: "upload files error",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				f.feeds.On("GetGeneration", a.ctx, a.generationID).
					Return(&entity.Generation{
						ID:            a.generationID,
						Type:          "test",
						Progress:      43,
						DataFetched:   true,
						FilesUploaded: 1,
						StartTime:     time.Unix(10, 0),
					}, nil)
				f.feeds.On("GetFactoryByGenerationType", "test").Return(f.factory, nil)
				f.feeds.On("UpdateGenerationState", a.ctx, &entity.Generation{
					ID:            a.generationID,
					Type:          "test",
					Progress:      0,
					DataFetched:   false,
					FilesUploaded: 0,
					StartTime:     time.Unix(10, 0),
				}).Return(nil)
				f.feeds.On("OnGenerationCanceled", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				f.factory.On("CreateDataFetcher", mock.Anything).Return(f.dataFetcher)
				f.factory.On("CreateFileFormatter", mock.Anything, mock.Anything).Return(f.fileFormatter)
				f.factory.On("CreateUploader", mock.Anything).Return(f.uploader)

				f.dataFetcher.
					On("StreamData", mock.Anything).Return(nil).
					On("OnDataFetched", mock.Anything).Return(nil).
					On("OnProgress", mock.Anything).Return(nil)
				f.fileFormatter.
					On("FormatFiles", mock.Anything).Return(nil)
				f.uploader.
					On("UploadFiles", mock.Anything).Return(defaultErr).After(time.Millisecond*5).
					On("OnUpload", mock.Anything).Return(nil)

				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)
			},
			wantErr: defaultErr,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fields := defaultFields()
			interactor := fields.newInteractor()
			testCase.setupMocks(testCase.args, fields)

			gotErr := interactor.RestartGeneration(testCase.args.ctx, testCase.args.generationID)

			assert.Equal(t, testCase.wantErr, gotErr)
			fields.assertExpectations(t)
		})
	}
}

func TestListGenerationTypes(t *testing.T) {
	testCases := []struct {
		name       string
		setupMocks func(*fields)
		want       interface{}
		wantErr    error
	}{
		{
			name: "succeed",
			setupMocks: func(f *fields) {
				f.feeds.On("ListAllowedTypes", mock.Anything).Return([]string{"a", "b"})
				f.presenter.On("PresentGenerationTypes", mock.Anything).
					Return(func(in []string) interface{} {
						return in
					})
			},
			want: []string{"a", "b"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fields := defaultFields()
			interactor := fields.newInteractor()
			testCase.setupMocks(fields)

			got, gotErr := interactor.ListGenerationTypes(context.Background())

			assert.Equal(t, testCase.want, got)
			assert.Equal(t, testCase.wantErr, gotErr)
			fields.assertExpectations(t)
		})
	}
}

func TestFeedInteractor_ListGenerations(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	defaultArgs := func() *args {
		return &args{ctx: context.Background()}
	}
	testCases := []struct {
		name       string
		args       *args
		setupMocks func(*args, *fields)
		want       interface{}
		wantErr    error
	}{
		{
			name: "succeed",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				f.feeds.
					On("ListGenerations", a.ctx).
					Return([]*entity.Generation{&generation1, &generation2}, nil)
				f.presenter.
					On("PresentListGenerations", mock.Anything).
					Return(func(out *interactor.ListGenerationsOut) interface{} {
						return out
					})
			},
			want: &interactor.ListGenerationsOut{
				{
					ID:        generation1.ID,
					Type:      generation1.Type,
					Progress:  generation1.Progress,
					StartTime: generation1.StartTime,
					EndTime:   generation1.EndTime,
				},
				{
					ID:        generation2.ID,
					Type:      generation2.Type,
					Progress:  generation2.Progress,
					StartTime: generation2.StartTime,
					EndTime:   generation2.EndTime,
				},
			},
		},
		{
			name: "feeds.ListGenerations error",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				f.feeds.On("ListGenerations", a.ctx).Return(nil, defaultErr)
				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)
			},
			wantErr: defaultErr,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fields := defaultFields()
			interactor := fields.newInteractor()
			testCase.setupMocks(testCase.args, fields)

			got, gotErr := interactor.ListGenerations(testCase.args.ctx)

			assert.Equal(t, testCase.want, got)
			assert.Equal(t, testCase.wantErr, gotErr)
			fields.assertExpectations(t)
		})
	}
}

func TestFeedInteractor_CancelGeneration(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	defaultArgs := func() *args {
		return &args{
			ctx: context.Background(),
			id:  defaultID,
		}
	}
	testCases := []struct {
		name       string
		args       *args
		setupMocks func(*args, *fields)
		wantErr    error
	}{
		{
			name: "succeed",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				f.feeds.On("CancelGeneration", a.ctx, a.id).Return(nil)
			},
		},
		{
			name: "feeds.CancelGeneration error",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				f.feeds.On("CancelGeneration", a.ctx, a.id).Return(defaultErr)
				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)
			},
			wantErr: defaultErr,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fields := defaultFields()
			interactor := fields.newInteractor()
			testCase.setupMocks(testCase.args, fields)

			gotErr := interactor.CancelGeneration(testCase.args.ctx, testCase.args.id)

			assert.Equal(t, testCase.wantErr, gotErr)
			fields.assertExpectations(t)
		})
	}
}
