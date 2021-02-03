package interactor_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-feedmaker/entity"
	helper "go-feedmaker/infrastructure/testing"
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
	uploader  *mocks.Uploader
	feeds     *mocks.FeedRepo
	presenter *mocks.Presenter
}

func defaultFields() *fields {
	return &fields{
		uploader:  new(mocks.Uploader),
		feeds:     new(mocks.FeedRepo),
		presenter: new(mocks.Presenter),
	}
}

func (f *fields) newInteractor() interactor.FeedInteractor {
	return interactor.NewFeedInteractor(f.uploader, f.feeds, f.presenter)
}

func (f *fields) assertExpectations(t *testing.T) {
	f.uploader.AssertExpectations(t)
	f.feeds.AssertExpectations(t)
	f.presenter.AssertExpectations(t)
}

func TestNewFeedInteractor(t *testing.T) {
	fields := defaultFields()
	i := interactor.NewFeedInteractor(fields.uploader, fields.feeds, fields.presenter)
	assert.Equal(t, fields.uploader, i.FileRepo())
	assert.Equal(t, fields.feeds, i.GenerationRepo())
	assert.Equal(t, fields.presenter, i.Presenter())
}

func TestFeedInteractor_Generate(t *testing.T) {
	t.SkipNow()
	type args struct {
		ctx            context.Context
		generationType string
	}
	defaultArgs := func() *args {
		return &args{ctx: context.Background(), generationType: "test"}
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
					return g.Progress == 0 && g.StartTime.Sub(time.Now()) < time.Second && g.Type == a.generationType
				}

				factory := new(mocks.FeedFactory)
				dataFetcher := new(mocks.DataFetcher)
				fileFormatter := new(mocks.FileFormatter)
				var dataStream <-chan []string
				fileStream := make(chan io.Reader, 2)

				file1 := helper.OpenFile(t, "testdata/eggs0.csv")
				file2 := helper.OpenFile(t, "testdata/eggs1.csv")
				fileStream <- file1
				fileStream <- file2
				close(fileStream)

				f.feeds.On("GetFactoryByGenerationType", a.generationType).Return(factory, nil)
				f.feeds.On("StoreGeneration", a.ctx, mock.MatchedBy(generationMatches)).Return(&generation1, nil)
				factory.On("CreateDataFetcher").Return(dataFetcher)
				factory.On("CreateFileFormatter", dataStream).Return(fileFormatter)
				dataFetcher.On("StreamData", a.ctx).Return(dataStream, nil)
				fileFormatter.On("FormatFiles", a.ctx).Return((<-chan io.Reader)(fileStream), nil)
				f.uploader.On("Upload", a.ctx, a.generationType, file1).Return(nil)
				f.uploader.On("Upload", a.ctx, a.generationType, file2).Return(nil)

				t.Cleanup(func() {
					factory.AssertExpectations(t)
					dataFetcher.AssertExpectations(t)
					fileFormatter.AssertExpectations(t)
				})
			},
		},
		{
			name: "feeds.StoreGeneration error",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				generationMatches := func(g *entity.Generation) bool {
					return g.Progress == 0 && g.StartTime.Sub(time.Now()) < time.Second && g.Type == a.generationType
				}
				factory := new(mocks.FeedFactory)
				f.feeds.On("GetFactoryByGenerationType", a.generationType).Return(factory, nil)
				f.feeds.On("StoreGeneration", a.ctx, mock.MatchedBy(generationMatches)).Return(nil, defaultErr)
				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)
				t.Cleanup(func() {
					factory.AssertExpectations(t)
				})
			},
			wantErr: defaultErr,
		},
		{
			name: "feeds.GetFactoryByGenerationType error",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				f.feeds.On("GetFactoryByGenerationType", a.generationType).Return(nil, defaultErr)
				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)
			},
			wantErr: defaultErr,
		},
		{
			name: "feeds.StreamData error",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				generationMatches := func(g *entity.Generation) bool {
					return g.Progress == 0 && g.StartTime.Sub(time.Now()) < time.Second && g.Type == a.generationType
				}
				factory := new(mocks.FeedFactory)
				dataFetcher := new(mocks.DataFetcher)
				fileFormatter := new(mocks.FileFormatter)

				f.feeds.On("GetFactoryByGenerationType", a.generationType).Return(factory, nil)
				f.feeds.On("StoreGeneration", a.ctx, mock.MatchedBy(generationMatches)).Return(&generation1, nil)
				factory.On("CreateDataFetcher").Return(dataFetcher)
				dataFetcher.On("StreamData", a.ctx).Return(nil, defaultErr)
				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)

				t.Cleanup(func() {
					factory.AssertExpectations(t)
					dataFetcher.AssertExpectations(t)
					fileFormatter.AssertExpectations(t)
				})
			},
			wantErr: defaultErr,
		},
		{
			name: "feeds.FormatFiles error",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				generationMatches := func(g *entity.Generation) bool {
					return g.Progress == 0 && g.StartTime.Sub(time.Now()) < time.Second && g.Type == a.generationType
				}
				factory := new(mocks.FeedFactory)
				dataFetcher := new(mocks.DataFetcher)
				fileFormatter := new(mocks.FileFormatter)
				var dataStream <-chan []string

				f.feeds.On("GetFactoryByGenerationType", a.generationType).Return(factory, nil)
				f.feeds.On("StoreGeneration", a.ctx, mock.MatchedBy(generationMatches)).Return(&generation1, nil)
				factory.On("CreateDataFetcher").Return(dataFetcher)
				factory.On("CreateFileFormatter", dataStream).Return(fileFormatter)
				dataFetcher.On("StreamData", a.ctx).Return(dataStream, nil)
				fileFormatter.On("FormatFiles", a.ctx).Return(nil, defaultErr)
				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)

				t.Cleanup(func() {
					factory.AssertExpectations(t)
					dataFetcher.AssertExpectations(t)
					fileFormatter.AssertExpectations(t)
				})
			},
			wantErr: defaultErr,
		},
		{
			name: "uploader.Upload error ignored",
			args: defaultArgs(),
			setupMocks: func(a *args, f *fields) {
				generationMatches := func(g *entity.Generation) bool {
					return g.Progress == 0 && g.StartTime.Sub(time.Now()) < time.Second && g.Type == a.generationType
				}
				factory := new(mocks.FeedFactory)
				dataFetcher := new(mocks.DataFetcher)
				fileFormatter := new(mocks.FileFormatter)
				var dataStream <-chan []string
				fileStream := make(chan io.Reader, 2)

				file1 := helper.OpenFile(t, "testdata/eggs0.csv")
				file2 := helper.OpenFile(t, "testdata/eggs1.csv")
				fileStream <- file1
				fileStream <- file2
				close(fileStream)

				f.feeds.On("GetFactoryByGenerationType", a.generationType).Return(factory, nil)
				f.feeds.On("StoreGeneration", a.ctx, mock.MatchedBy(generationMatches)).Return(&generation1, nil)
				factory.On("CreateDataFetcher").Return(dataFetcher)
				factory.On("CreateFileFormatter", dataStream).Return(fileFormatter)
				dataFetcher.On("StreamData", a.ctx).Return(dataStream, nil)
				fileFormatter.On("FormatFiles", a.ctx).Return((<-chan io.Reader)(fileStream), nil)
				f.uploader.On("Upload", a.ctx, a.generationType, file1).Return(defaultErr)
				f.uploader.On("Upload", a.ctx, a.generationType, file2).Return(nil)

				t.Cleanup(func() {
					factory.AssertExpectations(t)
					dataFetcher.AssertExpectations(t)
					fileFormatter.AssertExpectations(t)
				})
			},
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
				f.feeds.
					On("ListAllowedTypes", mock.Anything).
					Return([]string{"a", "b"}, nil)
				f.presenter.
					On("PresentGenerationTypes", mock.Anything).
					Return(func(in []string) interface{} {
						return in
					})
			},
			want: []string{"a", "b"},
		},
		{
			name: "feeds.ListAllowedTypes error",
			setupMocks: func(f *fields) {
				f.feeds.On("ListAllowedTypes", mock.Anything).Return(nil, defaultErr)
				f.presenter.On("PresentErr", mock.Anything).Return(errPassThrough)
			},
			wantErr: defaultErr,
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
