package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-feedmaker/infrastructure/rest"
	"go-feedmaker/infrastructure/scheduler"
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
	h := rest.NewHandler(fields.feeds, fields.scheduler)
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
			h := rest.NewHandler(testCase.fields.feeds, testCase.fields.scheduler)
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
			h := rest.NewHandler(testCase.fields.feeds, testCase.fields.scheduler)
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
		wantBody       []byte
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
			wantStatusCode: http.StatusCreated,
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
			wantBody:       mustMarshal(map[string]string{"details": defaultTestErr.Error()}),
		},
		{
			name:           "empty generation type",
			fields:         defaultHandlerFields(),
			setupMocks:     func(fields *handlerFields, args *args) {},
			args:           defaultArgs(""),
			wantStatusCode: http.StatusBadRequest,
			wantBody: mustMarshal(map[string]string{
				"details": fmt.Errorf("looking for generation-type: %w",
					rest.ErrValueNotFoundInURL).Error(),
			}),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields, testCase.args)
			h := rest.NewHandler(testCase.fields.feeds, testCase.fields.scheduler)
			h.GenerateFeed(testCase.args.w, testCase.args.r)
			gotStatusCode := testCase.args.w.Code
			gotBody := testCase.args.w.Body.Bytes()
			assert.Equal(t, testCase.wantStatusCode, gotStatusCode)
			assert.Equal(t, testCase.wantBody, gotBody)
			testCase.fields.feeds.AssertExpectations(t)
		})
	}
}

func Test_handler_CancelGeneration(t *testing.T) {
	type args struct {
		w            *httptest.ResponseRecorder
		r            *http.Request
		generationID string
	}
	defaultArgs := func(generationID string) *args {
		request := httptest.NewRequest(http.MethodDelete, "/generations/"+generationID, nil)
		if generationID != "" {
			vars := map[string]string{"generation-id": generationID}
			request = mux.SetURLVars(request, vars)
		}
		return &args{
			w:            httptest.NewRecorder(),
			r:            request,
			generationID: generationID,
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
					On("CancelGeneration", args.r.Context(), args.generationID).
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
					On("CancelGeneration", args.r.Context(), args.generationID).
					Return(defaultTestErr)
			},
			args:           defaultArgs("foobar"),
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       mustMarshal(map[string]string{"details": defaultTestErr.Error()}),
		},
		{
			name:           "empty generation id",
			fields:         defaultHandlerFields(),
			setupMocks:     func(fields *handlerFields, args *args) {},
			args:           defaultArgs(""),
			wantStatusCode: http.StatusBadRequest,
			wantBody: mustMarshal(map[string]string{
				"details": fmt.Errorf("looking for generation-id: %w",
					rest.ErrValueNotFoundInURL).Error(),
			}),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields, testCase.args)
			h := rest.NewHandler(testCase.fields.feeds, testCase.fields.scheduler)
			h.CancelGeneration(testCase.args.w, testCase.args.r)
			gotStatusCode := testCase.args.w.Code
			gotBody := testCase.args.w.Body.Bytes()
			assert.Equal(t, testCase.wantStatusCode, gotStatusCode)
			assert.Equal(t, testCase.wantBody, gotBody)
			testCase.fields.feeds.AssertExpectations(t)
		})
	}
}

func Test_handler_ScheduleGeneration(t *testing.T) {
	type args struct {
		w              *httptest.ResponseRecorder
		r              *http.Request
		generationType string
		scheduleIn     *rest.ScheduleTaskIn
		ctx            context.Context
	}
	defaultArgs := func(generationType string, scheduleIn *rest.ScheduleTaskIn) *args {
		var bodyContent []byte
		if scheduleIn != nil {
			bodyContent = mustMarshal(scheduleIn)
		}
		body := bytes.NewBuffer(bodyContent)
		request := httptest.NewRequest(http.MethodPost, "/generations/"+generationType+"/schedule", body)
		if generationType != "" {
			vars := map[string]string{"generation-type": generationType}
			request = mux.SetURLVars(request, vars)
		}
		return &args{
			w:              httptest.NewRecorder(),
			r:              request,
			generationType: generationType,
			scheduleIn:     scheduleIn,
			ctx:            context.Background(),
		}
	}
	testCases := []struct {
		name           string
		fields         *handlerFields
		args           *args
		setupMocks     func(*handlerFields, *args)
		wantStatusCode int
		wantBody       []byte
	}{
		{
			name:   "succeed",
			fields: defaultHandlerFields(),
			args:   defaultArgs("foobar", &defaultScheduleIn),
			setupMocks: func(fields *handlerFields, args *args) {
				fields.scheduler.
					On("ScheduleTask", scheduler.TaskID(args.generationType), mock.Anything).
					Return(nil)
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name:   "error in scheduler.ScheduleTask",
			fields: defaultHandlerFields(),
			args:   defaultArgs("foobar", &defaultScheduleIn),
			setupMocks: func(fields *handlerFields, args *args) {
				fields.scheduler.
					On("ScheduleTask", scheduler.TaskID(args.generationType), mock.Anything).
					Return(defaultTestErr)
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       mustMarshal(map[string]string{"details": defaultTestErr.Error()}),
		},
		{
			name:           "empty generation type",
			fields:         defaultHandlerFields(),
			setupMocks:     func(fields *handlerFields, args *args) {},
			args:           defaultArgs("", &defaultScheduleIn),
			wantStatusCode: http.StatusBadRequest,
			wantBody: mustMarshal(map[string]string{
				"details": fmt.Errorf("looking for generation-type: %w",
					rest.ErrValueNotFoundInURL).Error(),
			}),
		},
		{
			name:           "empty request body",
			fields:         defaultHandlerFields(),
			setupMocks:     func(fields *handlerFields, args *args) {},
			args:           defaultArgs("foobar", nil),
			wantStatusCode: http.StatusBadRequest,
			wantBody: mustMarshal(map[string]string{
				"details": fmt.Errorf("%w: EOF",
					rest.ErrReadingRequestBody).Error(),
			}),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields, testCase.args)
			h := rest.NewHandler(testCase.fields.feeds, testCase.fields.scheduler)
			h.ScheduleGeneration(testCase.args.w, testCase.args.r)
			gotStatusCode := testCase.args.w.Code
			gotBody := testCase.args.w.Body.Bytes()
			assert.Equal(t, testCase.wantStatusCode, gotStatusCode)
			assert.Equal(t, testCase.wantBody, gotBody)
			testCase.fields.feeds.AssertExpectations(t)
		})
	}
}

func Test_handler_ListSchedules(t *testing.T) {
	type args struct {
		w              *httptest.ResponseRecorder
		r              *http.Request
		generationType string
		scheduleIn     *rest.ScheduleTaskIn
	}
	defaultArgs := func() *args {
		return &args{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest(http.MethodPost, "/generations/schedules", nil),
		}
	}
	testCases := []struct {
		name           string
		fields         *handlerFields
		args           *args
		setupMocks     func(*handlerFields, *args)
		wantStatusCode int
		wantBody       []byte
	}{
		{
			name:   "succeed",
			fields: defaultHandlerFields(),
			args:   defaultArgs(),
			setupMocks: func(fields *handlerFields, args *args) {
				fields.scheduler.
					On("ListSchedules").
					Return(defaultTaskSchedules, nil)
			},
			wantStatusCode: http.StatusCreated,
			wantBody:       mustMarshal(rest.MakeSchedulesOut(defaultTaskSchedules)),
		},
		{
			name:   "error in scheduler.ListSchedules",
			fields: defaultHandlerFields(),
			args:   defaultArgs(),
			setupMocks: func(fields *handlerFields, args *args) {
				fields.scheduler.
					On("ListSchedules").
					Return(nil, defaultTestErr)
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       mustMarshal(map[string]string{"details": defaultTestErr.Error()}),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields, testCase.args)
			h := rest.NewHandler(testCase.fields.feeds, testCase.fields.scheduler)
			h.ListSchedules(testCase.args.w, testCase.args.r)
			gotStatusCode := testCase.args.w.Code
			gotBody := testCase.args.w.Body.Bytes()
			assert.Equal(t, testCase.wantStatusCode, gotStatusCode)
			assert.JSONEq(t, string(testCase.wantBody), string(gotBody))
			testCase.fields.feeds.AssertExpectations(t)
		})
	}
}

func Test_handler_UnscheduleGeneration(t *testing.T) {
	type args struct {
		w              *httptest.ResponseRecorder
		r              *http.Request
		generationType string
	}
	defaultArgs := func(generationType string) *args {
		request := httptest.NewRequest(http.MethodPost, "/generations/"+generationType+"/schedule", nil)
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
		args           *args
		setupMocks     func(*handlerFields, *args)
		wantStatusCode int
		wantBody       []byte
	}{
		{
			name:   "succeed",
			fields: defaultHandlerFields(),
			args:   defaultArgs("foobar"),
			setupMocks: func(fields *handlerFields, args *args) {
				fields.scheduler.
					On("RemoveTask", scheduler.TaskID(args.generationType)).
					Return(nil)
			},
			wantStatusCode: http.StatusAccepted,
		},
		{
			name:   "error in scheduler.RemoveTask",
			fields: defaultHandlerFields(),
			args:   defaultArgs("foobar"),
			setupMocks: func(fields *handlerFields, args *args) {
				fields.scheduler.
					On("RemoveTask", scheduler.TaskID(args.generationType)).
					Return(defaultTestErr)
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       mustMarshal(map[string]string{"details": defaultTestErr.Error()}),
		},
		{
			name:           "empty generation type",
			fields:         defaultHandlerFields(),
			setupMocks:     func(fields *handlerFields, args *args) {},
			args:           defaultArgs(""),
			wantStatusCode: http.StatusBadRequest,
			wantBody: mustMarshal(map[string]string{
				"details": fmt.Errorf("looking for generation-type: %w",
					rest.ErrValueNotFoundInURL).Error(),
			}),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields, testCase.args)
			h := rest.NewHandler(testCase.fields.feeds, testCase.fields.scheduler)
			h.UnscheduleGeneration(testCase.args.w, testCase.args.r)
			gotStatusCode := testCase.args.w.Code
			gotBody := testCase.args.w.Body.Bytes()
			assert.Equal(t, testCase.wantStatusCode, gotStatusCode)
			assert.Equal(t, testCase.wantBody, gotBody)
			testCase.fields.feeds.AssertExpectations(t)
		})
	}
}
