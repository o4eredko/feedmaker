package scheduler_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-feedmaker/infrastructure/scheduler"
	"go-feedmaker/infrastructure/scheduler/mocks"
)

type (
	cmdFields struct {
		f    *mocks.CmdFunc
		args []interface{}
	}
)

var (
	defaultContext = context.Background()
	defaultArgs    = []interface{}{defaultContext, "foobar", 42}
)

func defaultCmdFields() *cmdFields {
	mockedCmdFunc := new(mocks.CmdFunc)
	return &cmdFields{
		f:    mockedCmdFunc,
		args: defaultArgs,
	}
}

func TestNewCmd(t *testing.T) {
	getArgs := func(arguments []reflect.Value) []interface{} {
		args := make([]interface{}, len(arguments))
		for i, arg := range arguments {
			args[i] = arg.Interface()
		}
		return args
	}
	testCases := []struct {
		name    string
		fields  *cmdFields
		wantErr error
	}{
		{
			name:   "succeed",
			fields: defaultCmdFields(),
		},
		{
			name: "arguments amount mismatch",
			fields: &cmdFields{
				f:    new(mocks.CmdFunc),
				args: []interface{}{"foo", 42},
			},
			wantErr: scheduler.ErrArgumentsAmountMismatch,
		},
		{
			name: "invalid argument type",
			fields: &cmdFields{
				f:    new(mocks.CmdFunc),
				args: []interface{}{defaultContext, "foo", "42"},
			},
			wantErr: scheduler.ErrInvalidArgumentType,
		},
		{
			name: "argument is not assignable to expected interface",
			fields: &cmdFields{
				f:    new(mocks.CmdFunc),
				args: []interface{}{"defaultContext", "foo", 42},
			},
			wantErr: scheduler.ErrInvalidArgumentType,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cmd, gotErr := scheduler.NewCmd(testCase.fields.f.Execute, testCase.fields.args...)
			assert.True(t, errors.Is(gotErr, testCase.wantErr), "want: %v\ngot: %v", testCase.wantErr, gotErr)
			if gotErr == nil {
				wantFunc := reflect.ValueOf(testCase.fields.f.Execute).Pointer()
				gotFunc := reflect.ValueOf(cmd.Func.Interface()).Pointer()
				assert.Equal(t, wantFunc, gotFunc)
				assert.Equal(t, testCase.fields.args, getArgs(cmd.Args))
			}
		})
	}
}

func TestCmd_Run(t *testing.T) {
	testCases := []struct {
		name       string
		fields     *cmdFields
		setupMocks func(fields *cmdFields)
	}{
		{
			name:   "succeed",
			fields: defaultCmdFields(),
			setupMocks: func(fields *cmdFields) {
				fields.f.On("Execute", fields.args...)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields)
			cmd, err := scheduler.NewCmd(testCase.fields.f.Execute, testCase.fields.args...)
			assert.NoError(t, err)
			cmd.Run()
			testCase.fields.f.AssertExpectations(t)
		})
	}
}
