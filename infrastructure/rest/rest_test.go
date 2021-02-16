package rest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-feedmaker/infrastructure/rest"
	"go-feedmaker/infrastructure/rest/mocks"
)

func TestNewAPIServer(t *testing.T) {
	config := &rest.Config{Host: "foo.bar.baz", Port: "4213"}
	handler := new(mocks.Handler)
	wsHandler := new(mocks.WSHandler)
	router := rest.NewRouter(handler, wsHandler)
	server := rest.NewAPIServer(config, router).Server()
	assert.Equal(t, config.Addr(), server.Addr)
	assert.Equal(t, router, server.Handler)
}
