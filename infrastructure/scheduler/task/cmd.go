package task

import (
	"errors"
	"fmt"
	"reflect"
)

type (
	Cmd struct {
		f    reflect.Value
		args []reflect.Value
	}
)

var (
	ErrArgumentsAmountMismatch = errors.New("arguments amount mismatch")
	ErrInvalidArgumentType     = errors.New("invalid argument type")
)

func NewCmd(function interface{}, arguments ...interface{}) (*Cmd, error) {
	args := make([]reflect.Value, len(arguments))
	for i, arg := range arguments {
		args[i] = reflect.ValueOf(arg)
	}
	cmd := &Cmd{f: reflect.ValueOf(function), args: args}
	if err := cmd.validate(); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (c *Cmd) Run() {
	c.f.Call(c.args)
}

func (c *Cmd) validate() error {
	funcType := reflect.TypeOf(c.f.Interface())
	numIn := funcType.NumIn()
	numArgs := len(c.args)
	if numIn != numArgs {
		return fmt.Errorf("%w: expected %d arguments but got %d",
			ErrArgumentsAmountMismatch, numIn, numArgs)
	}
	for i, arg := range c.args {
		wantArgType := funcType.In(i).Kind()
		gotArgType := arg.Kind()
		if wantArgType != gotArgType {
			return fmt.Errorf("%w: at position %d: expected %v but got %v",
				ErrInvalidArgumentType, i, wantArgType, gotArgType)
		}
	}
	return nil
}
