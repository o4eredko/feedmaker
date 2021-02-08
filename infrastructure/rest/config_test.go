package rest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-feedmaker/infrastructure/rest"
)

func TestApi_Addr(t *testing.T) {
	config := rest.Config{Host: "localhost", Port: "8888"}
	assert.Equal(t, "localhost:8888", config.Addr())
}

func TestApi_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		obj     rest.Config
		wantErr bool
	}{
		{
			name: "succeed",
			obj:  rest.Config{Host: "foobar", Port: "42"},
		},
		{
			name:    "without host",
			obj:     rest.Config{Port: "42"},
			wantErr: true,
		},
		{
			name:    "without port",
			obj:     rest.Config{Host: "foobar"},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.obj.Validate()
			assert.Equal(t, got != nil, tc.wantErr, got)
		})
	}
}
