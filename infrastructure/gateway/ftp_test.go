package gateway_test

import (
	"bytes"
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

type fields struct {
	dialer     *mocks.Dialer
	config     gateway.FtpConfig
	connection *mocks.FtpConnection
}

func defaultFields() *fields {
	return &fields{
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
		fields     *fields
		setupMocks func(f *fields)
		wantErr    error
	}{
		{
			name:   "succeed",
			fields: defaultFields(),
			setupMocks: func(f *fields) {
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
			fields: defaultFields(),
			setupMocks: func(f *fields) {
				f.dialer.
					On("DialTimeout", f.config.Addr(), f.config.ConnTimeout).
					Return(nil, defaultErr)
			},
			wantErr: defaultErr,
		},
		{
			name:   "login error",
			fields: defaultFields(),
			setupMocks: func(f *fields) {
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
		path string
		r    io.Reader
	}
	defaultArgs := func() *args {
		return &args{
			path: "test",
			r:    bytes.NewBufferString("test"),
		}
	}
	testCases := []struct {
		name       string
		args       *args
		fields     *fields
		setupMocks func(*args, *fields)
		wantErr    error
	}{
		{
			name: "succeed",
			args: defaultArgs(),
			fields: &fields{
				config:     gateway.FtpConfig{},
				dialer:     new(mocks.Dialer),
				connection: new(mocks.FtpConnection),
			},
			setupMocks: func(a *args, f *fields) {
				f.connection.On("Stor", a.path, a.r).Return(nil)
			},
		},
		{
			name: "disconnected error",
			args: defaultArgs(),
			fields: &fields{
				config: gateway.FtpConfig{},
				dialer: new(mocks.Dialer),
			},
			setupMocks: func(a *args, f *fields) {},
			wantErr:    gateway.ErrFtpDisconnected,
		},
		{
			name: "Stor error",
			args: defaultArgs(),
			fields: &fields{
				config:     gateway.FtpConfig{},
				dialer:     new(mocks.Dialer),
				connection: new(mocks.FtpConnection),
			},
			setupMocks: func(a *args, f *fields) {
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

			gotErr := ftpGateway.Upload(testCase.args.path, testCase.args.r)

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
		fields     *fields
		setupMocks func(f *fields)
		wantErr    error
	}{
		{
			name: "succeed",
			fields: &fields{
				config:     gateway.FtpConfig{},
				dialer:     new(mocks.Dialer),
				connection: new(mocks.FtpConnection),
			},
			setupMocks: func(f *fields) {
				f.connection.On("Quit").Return(nil)
			},
		},
		{
			name: "disconnected error",
			fields: &fields{
				config:     gateway.FtpConfig{},
				dialer:     new(mocks.Dialer),
				connection: nil,
			},
			setupMocks: func(f *fields) {},
			wantErr:    gateway.ErrFtpDisconnected,
		},
		{
			name: "Quit error",
			fields: &fields{
				config:     gateway.FtpConfig{},
				dialer:     new(mocks.Dialer),
				connection: new(mocks.FtpConnection),
			},
			setupMocks: func(f *fields) {
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
