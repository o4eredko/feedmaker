package interactor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-feedmaker/interactor"
	"go-feedmaker/interactor/mocks"
)

var (
	defaultErr     = errors.New("default test error")
	errPassThrough = func(err error) error {
		return err
	}
)

type fields struct {
	csvFetcher  *mocks.CsvRepo
	files       *mocks.FileRepo
	generations *mocks.GenerationRepo
	presenter   *mocks.Presenter
}

func defaultFields() *fields {
	return &fields{
		csvFetcher:  new(mocks.CsvRepo),
		files:       new(mocks.FileRepo),
		generations: new(mocks.GenerationRepo),
		presenter:   new(mocks.Presenter),
	}
}

func (f *fields) newInteractor() interactor.FeedInteractor {
	return interactor.NewFeedInteractor(f.csvFetcher, f.files, f.generations, f.presenter)
}

func (f *fields) assertExpectations(t *testing.T) {
	f.csvFetcher.AssertExpectations(t)
	f.files.AssertExpectations(t)
	f.generations.AssertExpectations(t)
	f.presenter.AssertExpectations(t)
}

func TestNewFeedInteractor(t *testing.T) {
	fields := defaultFields()
	i := interactor.NewFeedInteractor(fields.csvFetcher, fields.files, fields.generations, fields.presenter)
	assert.Equal(t, fields.csvFetcher, i.CsvFetcher())
	assert.Equal(t, fields.files, i.FileRepo())
	assert.Equal(t, fields.generations, i.GenerationRepo())
	assert.Equal(t, fields.presenter, i.Presenter())
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
				f.generations.
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
			name: "generations.ListAllowedTypes error",
			setupMocks: func(f *fields) {
				f.generations.On("ListAllowedTypes", mock.Anything).Return(nil, defaultErr)
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
