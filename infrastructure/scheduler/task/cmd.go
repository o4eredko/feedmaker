package task

import (
	"errors"
	"fmt"
	"reflect"
)

type (
	Cmd struct {
		Func reflect.Value
		Args []reflect.Value
	}
)

var (
	ErrArgumentsAmountMismatch = errors.New("arguments amount mismatch")
	ErrInvalidArgumentType     = errors.New("invalid argument type")
)

func NewCmd(function interface{}, arguments ...interface{}) (*Cmd, error) {
	f := makeFunc(function)
	args := makeArgs(arguments)
	if err := validate(f, args); err != nil {
		return nil, err
	}
	cmd := &Cmd{Func: f, Args: args}
	return cmd, nil
}

func (c *Cmd) Run() {
	c.Func.Call(c.Args)
}

func makeFunc(function interface{}) reflect.Value {
	return reflect.ValueOf(function)
}

func makeArgs(arguments []interface{}) []reflect.Value {
	args := make([]reflect.Value, len(arguments))
	for i, arg := range arguments {
		args[i] = reflect.ValueOf(arg)
	}
	return args
}

func validate(function reflect.Value, arguments []reflect.Value) error {
	funcType := reflect.TypeOf(function.Interface())
	numIn := funcType.NumIn()
	numArgs := len(arguments)
	if numIn != numArgs {
		return fmt.Errorf("%w: expected %d arguments but got %d",
			ErrArgumentsAmountMismatch, numIn, numArgs)
	}
	for i, arg := range arguments {
		wantArgType := funcType.In(i).Kind()
		gotArgType := arg.Kind()
		if wantArgType != gotArgType {
			return fmt.Errorf("%w: at position %d: expected %v but got %v",
				ErrInvalidArgumentType, i, wantArgType, gotArgType)
		}
	}
	return nil
}
