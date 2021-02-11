package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type CmdFunc struct {
	mock.Mock
}

func (_m *CmdFunc) Execute(ctx context.Context, stringArg string, intArg int) {
	_m.Called(ctx, stringArg, intArg)
}
