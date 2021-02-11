package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"go-feedmaker/infrastructure/rest"
	"go-feedmaker/infrastructure/rest/mocks"
)

func TestNewRouter(t *testing.T) {
	type args struct {
		handler *mocks.Handler
		request *http.Request
	}
	mustMakeArgs := func(method, path string) *args {
		request, err := http.NewRequest(method, path, nil)
		if err != nil {
			panic(err)
		}
		return &args{
			handler: new(mocks.Handler),
			request: request,
		}
	}
	testCases := []struct {
		name       string
		args       *args
		setupMocks func(*args)
	}{
		{
			name: "GET /generations",
			args: mustMakeArgs(http.MethodGet, "/generations"),
			setupMocks: func(args *args) {
				args.handler.On("ListGenerations", mock.Anything, mock.Anything)
			},
		},
		{
			name: "GET /generations/types",
			args: mustMakeArgs(http.MethodGet, "/generations/types"),
			setupMocks: func(args *args) {
				args.handler.On("ListGenerationTypes", mock.Anything, mock.Anything)
			},
		},
		{
			name: "POST /generations/foobar",
			args: mustMakeArgs(http.MethodPost, "/generations/foobar"),
			setupMocks: func(args *args) {
				args.handler.On("GenerateFeed", mock.Anything, mock.Anything)
			},
		},
		{
			name: "DELETE /generations/foobar",
			args: mustMakeArgs(http.MethodDelete, "/generations/foobar"),
			setupMocks: func(args *args) {
				args.handler.On("CancelGeneration", mock.Anything, mock.Anything)
			},
		},
		{
			name: "POST /generations/foobar/schedules",
			args: mustMakeArgs(http.MethodPost, "/generations/foobar/schedules"),
			setupMocks: func(args *args) {
				args.handler.On("ScheduleGeneration", mock.Anything, mock.Anything)
			},
		},
		{
			name: "GET /generations/schedules",
			args: mustMakeArgs(http.MethodGet, "/generations/schedules"),
			setupMocks: func(args *args) {
				args.handler.On("ListSchedules", mock.Anything, mock.Anything)
			},
		},
		{
			name: "DELETE /generations/foobar/schedules",
			args: mustMakeArgs(http.MethodDelete, "/generations/foobar/schedules"),
			setupMocks: func(args *args) {
				args.handler.On("UnscheduleGeneration", mock.Anything, mock.Anything)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.args)
			router := rest.NewRouter(testCase.args.handler)

			router.ServeHTTP(httptest.NewRecorder(), testCase.args.request)

			testCase.args.handler.AssertExpectations(t)
		})
	}
}
