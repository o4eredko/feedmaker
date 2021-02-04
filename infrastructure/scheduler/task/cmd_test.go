package task_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-feedmaker/infrastructure/scheduler/mocks"
	"go-feedmaker/infrastructure/scheduler/task"
)

type (
	cmdFields struct {
		f    *mocks.CmdFunc
		args []interface{}
	}
)

func defaultCmdFields() *cmdFields {
	mockedCmdFunc := new(mocks.CmdFunc)
	return &cmdFields{
		f:    mockedCmdFunc,
		args: []interface{}{"foo", "bar", "baz"},
	}
}

func TestNewCmd(t *testing.T) {
	type fields struct {
		f    interface{}
		args []interface{}
	}
	testCases := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "succeed",
			fields: fields{
				f:    func(string, int) {},
				args: []interface{}{"foo", 42},
			},
		},
		{
			name: "arguments amount mismatch",
			fields: fields{
				f:    func(string, int) {},
				args: []interface{}{"foo", 42, 13},
			},
			wantErr: task.ErrArgumentsAmountMismatch,
		},
		{
			name: "invalid argument type",
			fields: fields{
				f:    func(string, int) {},
				args: []interface{}{"foo", "42"},
			},
			wantErr: task.ErrInvalidArgumentType,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cmd, gotErr := task.NewCmd(testCase.fields.f, testCase.fields.args...)
			assert.Equal(t, testCase.wantErr, errors.Unwrap(gotErr))
			if gotErr == nil {
				wantFunc := reflect.ValueOf(testCase.fields.f).Pointer()
				gotFunc := reflect.ValueOf(cmd.CmdFunc()).Pointer()
				assert.Equal(t, wantFunc, gotFunc)
				assert.Equal(t, testCase.fields.args, cmd.Args())
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
				fields.f.On("Execute", fields.args)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setupMocks(testCase.fields)
			cmd, _ := task.NewCmd(testCase.fields.f.Execute, testCase.fields.args)
			cmd.Run()

			testCase.fields.f.AssertExpectations(t)
		})
	}
}
