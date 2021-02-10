package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"go-feedmaker/infrastructure/rest"
)

func mustMarshal(v interface{}) []byte {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func TestNewHandler(t *testing.T) {
	fields := defaultHandlerFields()
	h := rest.NewHandler(fields.feeds)
	assert.Equal(t, fields.feeds, h.Feeds())
}

func Test_handler_ListGenerations(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	defaultArgs := func() *args {
		return &args{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest(http.MethodGet, "/generations", nil),
		}
	}
	testCases := []struct {
		name           string
		fields         *handlerFields
		setupMocks     func(*handlerFields, *args)
		args           *args
		wantStatusCode int
		wantBody       []byte
	}{
		{
			name:   "succeed",
			fields: defaultHandlerFields(),
			setupMocks: func(fields *handlerFields, args *args) {
				fields.feeds.
					On("ListGenerations", args.r.Context()).
					Return(defaultSentinel, nil)
			},
			args:           defaultArgs(),
			wantStatusCode: http.StatusOK,
			wantBody:       mustMarshal(defaultSentinel),
		},
		{
			name:   "error in feeds.ListGenerations",
			fields: defaultHandlerFields(),
			setupMocks: func(fields *handlerFields, args *args) {
				fields.feeds.
					On("ListGenerations", args.r.Context()).
					Return(nil, defaultTestErr)
			},
			args:           defaultArgs(),
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       mustMarshal(map[string]string{"details": defaultTestErr.Error()}),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields, testCase.args)
			h := rest.NewHandler(testCase.fields.feeds)
			h.ListGenerations(testCase.args.w, testCase.args.r)
			gotStatusCode := testCase.args.w.Code
			gotBody := testCase.args.w.Body.Bytes()
			assert.Equal(t, testCase.wantStatusCode, gotStatusCode)
			assert.Equal(t, testCase.wantBody, gotBody)
			testCase.fields.feeds.AssertExpectations(t)
		})
	}
}

func Test_handler_ListGenerationTypes(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	defaultArgs := func() *args {
		return &args{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest(http.MethodGet, "/generations/types", nil),
		}
	}
	testCases := []struct {
		name           string
		fields         *handlerFields
		setupMocks     func(*handlerFields, *args)
		args           *args
		wantStatusCode int
		wantBody       []byte
	}{
		{
			name:   "succeed",
			fields: defaultHandlerFields(),
			setupMocks: func(fields *handlerFields, args *args) {
				fields.feeds.
					On("ListGenerationTypes", args.r.Context()).
					Return(defaultSentinel, nil)
			},
			args:           defaultArgs(),
			wantStatusCode: http.StatusOK,
			wantBody:       mustMarshal(defaultSentinel),
		},
		{
			name:   "error in feeds.ListGenerationTypes",
			fields: defaultHandlerFields(),
			setupMocks: func(fields *handlerFields, args *args) {
				fields.feeds.
					On("ListGenerationTypes", args.r.Context()).
					Return(nil, defaultTestErr)
			},
			args:           defaultArgs(),
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       mustMarshal(map[string]string{"details": defaultTestErr.Error()}),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields, testCase.args)
			h := rest.NewHandler(testCase.fields.feeds)
			h.ListGenerationTypes(testCase.args.w, testCase.args.r)
			gotStatusCode := testCase.args.w.Code
			gotBody := testCase.args.w.Body.Bytes()
			assert.Equal(t, testCase.wantStatusCode, gotStatusCode)
			assert.Equal(t, testCase.wantBody, gotBody)
			testCase.fields.feeds.AssertExpectations(t)
		})
	}
}

func Test_handler_GenerateFeed(t *testing.T) {
	type args struct {
		w              *httptest.ResponseRecorder
		r              *http.Request
		generationType string
	}
	defaultArgs := func(generationType string) *args {
		request := httptest.NewRequest(http.MethodPost, "/generations/"+generationType, nil)
		if generationType != "" {
			vars := map[string]string{"generation-type": generationType}
			request = mux.SetURLVars(request, vars)
		}
		return &args{
			w:              httptest.NewRecorder(),
			r:              request,
			generationType: generationType,
		}
	}
	testCases := []struct {
		name           string
		fields         *handlerFields
		setupMocks     func(*handlerFields, *args)
		args           *args
		wantStatusCode int
	}{
		{
			name:   "succeed",
			fields: defaultHandlerFields(),
			setupMocks: func(fields *handlerFields, args *args) {
				fields.feeds.
					On("GenerateFeed", args.r.Context(), args.generationType).
					Return(nil)
			},
			args:           defaultArgs("foobar"),
			wantStatusCode: http.StatusOK,
		},
		{
			name:   "error in feeds.ListGenerationTypes",
			fields: defaultHandlerFields(),
			setupMocks: func(fields *handlerFields, args *args) {
				fields.feeds.
					On("GenerateFeed", args.r.Context(), args.generationType).
					Return(defaultTestErr)
			},
			args:           defaultArgs("foobar"),
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "empty generation type",
			fields:         defaultHandlerFields(),
			setupMocks:     func(fields *handlerFields, args *args) {},
			args:           defaultArgs(""),
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields, testCase.args)
			h := rest.NewHandler(testCase.fields.feeds)
			h.GenerateFeed(testCase.args.w, testCase.args.r)
			gotStatusCode := testCase.args.w.Code
			assert.Equal(t, testCase.wantStatusCode, gotStatusCode)
			testCase.fields.feeds.AssertExpectations(t)
		})
	}
}

func Test_handler_CancelGeneration(t *testing.T) {
	type args struct {
		w              *httptest.ResponseRecorder
		r              *http.Request
		generationType string
	}
	defaultArgs := func(generationID string) *args {
		request := httptest.NewRequest(http.MethodDelete, "/generations/"+generationID, nil)
		if generationID != "" {
			vars := map[string]string{"generation-id": generationID}
			request = mux.SetURLVars(request, vars)
		}
		return &args{
			w:              httptest.NewRecorder(),
			r:              request,
			generationType: generationID,
		}
	}
	testCases := []struct {
		name           string
		fields         *handlerFields
		setupMocks     func(*handlerFields, *args)
		args           *args
		wantStatusCode int
	}{
		{
			name:   "succeed",
			fields: defaultHandlerFields(),
			setupMocks: func(fields *handlerFields, args *args) {
				fields.feeds.
					On("CancelGeneration", args.r.Context(), args.generationType).
					Return(nil)
			},
			args:           defaultArgs("foobar"),
			wantStatusCode: http.StatusOK,
		},
		{
			name:   "error in feeds.ListGenerationTypes",
			fields: defaultHandlerFields(),
			setupMocks: func(fields *handlerFields, args *args) {
				fields.feeds.
					On("CancelGeneration", args.r.Context(), args.generationType).
					Return(defaultTestErr)
			},
			args:           defaultArgs("foobar"),
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "empty generation id",
			fields:         defaultHandlerFields(),
			setupMocks:     func(fields *handlerFields, args *args) {},
			args:           defaultArgs(""),
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields, testCase.args)
			h := rest.NewHandler(testCase.fields.feeds)
			h.CancelGeneration(testCase.args.w, testCase.args.r)
			gotStatusCode := testCase.args.w.Code
			assert.Equal(t, testCase.wantStatusCode, gotStatusCode)
			testCase.fields.feeds.AssertExpectations(t)
		})
	}
}
