package repository

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	ErrInvalidRecord = errors.New("invalid record")
)

type (
	PolicyCheckerConfig struct {
		URL string
	}

	Requester interface {
		Do(req *http.Request) (*http.Response, error)
	}

	PolicyChecker struct {
		Requester Requester
		Config    PolicyCheckerConfig
	}
)

func (p *PolicyChecker) ValidateRecord(ctx context.Context, record []string) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.Config.URL,
		bytes.NewBufferString(strings.Join(record, " | ")),
	)
	if err != nil {
		return err
	}
	response, err := p.Requester.Do(request)
	if err != nil {
		return err
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	if response.StatusCode != http.StatusOK {
		respBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s, %w", respBody, ErrInvalidRecord)
	}
	return nil
}
