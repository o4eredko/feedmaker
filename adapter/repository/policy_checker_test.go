package repository_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-feedmaker/adapter/repository"
	"go-feedmaker/adapter/repository/mocks"
)

type policyCheckerFields struct {
	config    repository.PolicyCheckerConfig
	requester *mocks.Requester
}

func defaultPolicyCheckerFields() *policyCheckerFields {
	return &policyCheckerFields{
		config:    repository.PolicyCheckerConfig{URL: "http://abc.com"},
		requester: new(mocks.Requester),
	}
}

func matchRequests(wantReq *http.Request) func(gotReq *http.Request) bool {
	return func(gotReq *http.Request) bool {
		return wantReq.Method == gotReq.Method
	}
}

func TestPolicyChecker_ValidateRecord(t *testing.T) {
	type args struct {
		ctx    context.Context
		record []string
	}
	testCases := []struct {
		name       string
		args       *args
		fields     *policyCheckerFields
		setupMocks func(*args, *policyCheckerFields)
		wantErr    error
	}{
		{
			name: "succeed",
			args: &args{
				ctx:    context.Background(),
				record: []string{"a", "b", "c"},
			},
			fields: defaultPolicyCheckerFields(),
			setupMocks: func(a *args, f *policyCheckerFields) {
				request, err := http.NewRequestWithContext(
					a.ctx, http.MethodPost,
					f.config.URL, bytes.NewBufferString("a | b | c"),
				)
				assert.NoError(t, err)
				response := &http.Response{StatusCode: http.StatusOK}

				f.requester.On("Do", mock.MatchedBy(matchRequests(request))).Return(response, nil)
			},
		},
		{
			name: "request error",
			args: &args{
				ctx:    context.Background(),
				record: []string{"a", "b", "c"},
			},
			fields: defaultPolicyCheckerFields(),
			setupMocks: func(a *args, f *policyCheckerFields) {
				request, err := http.NewRequestWithContext(
					a.ctx, http.MethodPost,
					f.config.URL, bytes.NewBufferString("a | b | c"),
				)
				assert.NoError(t, err)

				f.requester.On("Do", mock.MatchedBy(matchRequests(request))).Return(nil, defaultErr)
			},
			wantErr: defaultErr,
		},
		{
			name: "bad status code",
			args: &args{
				ctx:    context.Background(),
				record: []string{"a", "b", "c"},
			},
			fields: defaultPolicyCheckerFields(),
			setupMocks: func(a *args, f *policyCheckerFields) {
				request, err := http.NewRequestWithContext(
					a.ctx, http.MethodPost,
					f.config.URL, bytes.NewBufferString("a | b | c"),
				)
				assert.NoError(t, err)
				response := &http.Response{
					Body:       ioutil.NopCloser(strings.NewReader("test")),
					StatusCode: http.StatusInternalServerError,
				}
				f.requester.On("Do", mock.MatchedBy(matchRequests(request))).Return(response, nil)
			},
			wantErr: repository.ErrInvalidRecord,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks(tc.args, tc.fields)
			checker := repository.PolicyChecker{
				Requester: tc.fields.requester,
				Config:    tc.fields.config,
			}

			gotErr := checker.ValidateRecord(tc.args.ctx, tc.args.record)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			tc.fields.requester.AssertExpectations(t)
		})
	}
}
