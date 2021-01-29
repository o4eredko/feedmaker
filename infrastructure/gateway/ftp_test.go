package gateway_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-feedmaker/infrastructure/gateway"
	"go-feedmaker/infrastructure/gateway/mocks"
)

var (
	defaultErr = errors.New("test error")
)

type ftpFields struct {
	dialer     *mocks.Dialer
	config     gateway.FtpConfig
	connection *mocks.FtpConnection
}

func defaultFtpFields() *ftpFields {
	return &ftpFields{
		dialer:     new(mocks.Dialer),
		connection: new(mocks.FtpConnection),
		config: gateway.FtpConfig{
			Host:        "localhost",
			Port:        "21",
			ConnTimeout: time.Millisecond,
			Username:    "test",
			Password:    "test,",
		},
	}
}

func TestFtpGateway_Connect(t *testing.T) {
	testCases := []struct {
		name       string
		fields     *ftpFields
		setupMocks func(f *ftpFields)
		wantErr    error
	}{
		{
			name:   "succeed",
			fields: defaultFtpFields(),
			setupMocks: func(f *ftpFields) {
				f.dialer.
					On("DialTimeout", f.config.Addr(), f.config.ConnTimeout).
					Return(f.connection, nil)
				f.connection.
					On("Login", f.config.Username, f.config.Password).
					Return(nil)
			},
		},
		{
			name:   "dial error",
			fields: defaultFtpFields(),
			setupMocks: func(f *ftpFields) {
				f.dialer.
					On("DialTimeout", f.config.Addr(), f.config.ConnTimeout).
					Return(nil, defaultErr)
			},
			wantErr: defaultErr,
		},
		{
			name:   "login error",
			fields: defaultFtpFields(),
			setupMocks: func(f *ftpFields) {
				f.dialer.
					On("DialTimeout", f.config.Addr(), f.config.ConnTimeout).
					Return(f.connection, nil)
				f.connection.
					On("Login", f.config.Username, f.config.Password).
					Return(defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields)
			ftpGateway := gateway.FtpGateway{
				Dialer: testCase.fields.dialer,
				Config: testCase.fields.config,
			}

			gotErr := ftpGateway.Connect()

			assert.Equal(t, testCase.wantErr, gotErr)
			testCase.fields.dialer.AssertExpectations(t)
			testCase.fields.connection.AssertExpectations(t)
		})
	}
}

func TestFtpGateway_Upload(t *testing.T) {
	type args struct {
		ctx  context.Context
		path string
		r    io.Reader
	}
	defaultArgs := func() *args {
		return &args{
			ctx:  context.Background(),
			path: "test",
			r:    bytes.NewBufferString("test"),
		}
	}
	testCases := []struct {
		name       string
		args       *args
		fields     *ftpFields
		setupMocks func(*args, *ftpFields)
		wantErr    error
	}{
		{
			name: "succeed",
			args: defaultArgs(),
			fields: &ftpFields{
				config:     gateway.FtpConfig{},
				dialer:     new(mocks.Dialer),
				connection: new(mocks.FtpConnection),
			},
			setupMocks: func(a *args, f *ftpFields) {
				f.connection.On("Stor", a.path, a.r).Return(nil)
			},
		},
		{
			name: "disconnected error",
			args: defaultArgs(),
			fields: &ftpFields{
				config: gateway.FtpConfig{},
				dialer: new(mocks.Dialer),
			},
			setupMocks: func(a *args, f *ftpFields) {},
			wantErr:    gateway.ErrFtpDisconnected,
		},
		{
			name: "Stor error",
			args: defaultArgs(),
			fields: &ftpFields{
				config:     gateway.FtpConfig{},
				dialer:     new(mocks.Dialer),
				connection: new(mocks.FtpConnection),
			},
			setupMocks: func(a *args, f *ftpFields) {
				f.connection.On("Stor", a.path, a.r).Return(defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.args, testCase.fields)
			ftpGateway := gateway.FtpGateway{
				Dialer: testCase.fields.dialer,
				Config: testCase.fields.config,
			}
			ftpGateway.SetConnection(testCase.fields.connection)

			gotErr := ftpGateway.Upload(testCase.args.ctx, testCase.args.path, testCase.args.r)

			assert.Equal(t, testCase.wantErr, gotErr)
			testCase.fields.dialer.AssertExpectations(t)
			if testCase.fields.connection != nil {
				testCase.fields.connection.AssertExpectations(t)
			}
		})
	}
}

func TestFtpGateway_Disconnect(t *testing.T) {
	testCases := []struct {
		name       string
		fields     *ftpFields
		setupMocks func(f *ftpFields)
		wantErr    error
	}{
		{
			name: "succeed",
			fields: &ftpFields{
				config:     gateway.FtpConfig{},
				dialer:     new(mocks.Dialer),
				connection: new(mocks.FtpConnection),
			},
			setupMocks: func(f *ftpFields) {
				f.connection.On("Quit").Return(nil)
			},
		},
		{
			name: "disconnected error",
			fields: &ftpFields{
				config:     gateway.FtpConfig{},
				dialer:     new(mocks.Dialer),
				connection: nil,
			},
			setupMocks: func(f *ftpFields) {},
			wantErr:    gateway.ErrFtpDisconnected,
		},
		{
			name: "Quit error",
			fields: &ftpFields{
				config:     gateway.FtpConfig{},
				dialer:     new(mocks.Dialer),
				connection: new(mocks.FtpConnection),
			},
			setupMocks: func(f *ftpFields) {
				f.connection.On("Quit").Return(defaultErr)
			},
			wantErr: defaultErr,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields)
			ftpGateway := gateway.FtpGateway{
				Dialer: testCase.fields.dialer,
				Config: testCase.fields.config,
			}
			ftpGateway.SetConnection(testCase.fields.connection)

			gotErr := ftpGateway.Disconnect()

			assert.Equal(t, testCase.wantErr, gotErr)
			testCase.fields.dialer.AssertExpectations(t)
			if testCase.fields.connection != nil {
				testCase.fields.connection.AssertExpectations(t)
			}
		})
	}
}
