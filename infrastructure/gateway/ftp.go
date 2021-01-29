package gateway

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"
)

type (
	FtpConfig struct {
		Host        string
		Port        string
		ConnTimeout time.Duration
		Username    string
		Password    string
	}

	Dialer interface {
		DialTimeout(addr string, timeout time.Duration) (FtpConnection, error)
	}

	FtpConnection interface {
		Login(user, password string) error
		Stor(path string, r io.Reader) error
		Quit() error
	}

	FtpGateway struct {
		Dialer     Dialer
		Config     FtpConfig
		connection FtpConnection
	}
)

var (
	ErrFtpDisconnected = errors.New("gateway is not connected to FTP")
)

func (c FtpConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (f *FtpGateway) Connect() error {
	connection, err := f.Dialer.DialTimeout(f.Config.Addr(), f.Config.ConnTimeout)
	if err != nil {
		return err
	}
	err = connection.Login(f.Config.Username, f.Config.Password)
	if err != nil {
		return err
	}
	f.connection = connection
	return nil
}

func (f *FtpGateway) Upload(ctx context.Context, path string, r io.Reader) error {
	if reflect.ValueOf(f.connection).IsNil() {
		return ErrFtpDisconnected
	}
	return f.connection.Stor(path, r)
}

func (f *FtpGateway) Disconnect() error {
	if reflect.ValueOf(f.connection).IsNil() {
		return ErrFtpDisconnected
	}
	return f.connection.Quit()
}
