package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"go-feedmaker/infrastructure/rest"
	"go-feedmaker/infrastructure/rest/mocks"
)

type (
	routerFields struct {
		handler *mocks.Handler
	}
)

func defaultRouterFields() *routerFields {
	return &routerFields{
		handler: new(mocks.Handler),
	}
}

func TestNewRouter(t *testing.T) {
	type args struct {
		responseWriter *httptest.ResponseRecorder
		request        *http.Request
	}
	mustMakeArgs := func(method, path string) *args {
		request, err := http.NewRequest(method, path, nil)
		if err != nil {
			panic(err)
		}
		return &args{
			responseWriter: httptest.NewRecorder(),
			request:        request,
		}
	}
	testCases := []struct {
		name       string
		fields     *routerFields
		args       *args
		setupMocks func(*routerFields)
	}{
		{
			name:   "GET /generations",
			fields: defaultRouterFields(),
			args:   mustMakeArgs(http.MethodGet, "/generations"),
			setupMocks: func(fields *routerFields) {
				fields.handler.On("ListGenerations", mock.Anything, mock.Anything)
			},
		},
		{
			name:   "GET /generations/types",
			fields: defaultRouterFields(),
			args:   mustMakeArgs(http.MethodGet, "/generations/types"),
			setupMocks: func(fields *routerFields) {
				fields.handler.On("ListGenerationTypes", mock.Anything, mock.Anything)
			},
		},
		{
			name:   "POST /generations/foobar",
			fields: defaultRouterFields(),
			args:   mustMakeArgs(http.MethodPost, "/generations/foobar"),
			setupMocks: func(fields *routerFields) {
				fields.handler.On("GenerateFeed", mock.Anything, mock.Anything)
			},
		},
		{
			name:   "DELETE /generations/foobar",
			fields: defaultRouterFields(),
			args:   mustMakeArgs(http.MethodDelete, "/generations/foobar"),
			setupMocks: func(fields *routerFields) {
				fields.handler.On("CancelGeneration", mock.Anything, mock.Anything)
			},
		},
		{
			name:   "POST /generations/foobar/schedules",
			fields: defaultRouterFields(),
			args:   mustMakeArgs(http.MethodPost, "/generations/foobar/schedules"),
			setupMocks: func(fields *routerFields) {
				fields.handler.On("ScheduleGeneration", mock.Anything, mock.Anything)
			},
		},
		{
			name:   "GET /generations/schedules",
			fields: defaultRouterFields(),
			args:   mustMakeArgs(http.MethodGet, "/generations/schedules"),
			setupMocks: func(fields *routerFields) {
				fields.handler.On("ListSchedules", mock.Anything, mock.Anything)
			},
		},
		{
			name:   "DELETE /generations/foobar/schedules",
			fields: defaultRouterFields(),
			args:   mustMakeArgs(http.MethodDelete, "/generations/foobar/schedules"),
			setupMocks: func(fields *routerFields) {
				fields.handler.On("UnscheduleGeneration", mock.Anything, mock.Anything)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields)
			router := rest.NewRouter(testCase.fields.handler)

			router.ServeHTTP(testCase.args.responseWriter, testCase.args.request)

			testCase.fields.handler.AssertExpectations(t)
		})
	}
}
