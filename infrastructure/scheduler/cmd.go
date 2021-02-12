package scheduler

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
	if err := checkNumberOfArgs(funcType, arguments); err != nil {
		return err
	}
	return checkTypesOfArgs(arguments, funcType)
}

func checkNumberOfArgs(funcType reflect.Type, arguments []reflect.Value) error {
	numIn := funcType.NumIn()
	numArgs := len(arguments)
	if numIn != numArgs {
		return fmt.Errorf("%w: expected %d arguments but got %d",
			ErrArgumentsAmountMismatch, numIn, numArgs)
	}
	return nil
}

func checkTypesOfArgs(arguments []reflect.Value, funcType reflect.Type) error {
	for i, arg := range arguments {
		wantArgType := funcType.In(i)
		gotArgType := arg.Type()
		if err := checkTypes(wantArgType, gotArgType); err != nil {
			return fmt.Errorf("at position %d: %w", i, err)
		}
	}
	return nil
}

func checkTypes(wantArgType, gotArgType reflect.Type) error {
	var assertion func(reflect.Type, reflect.Type) error
	switch wantArgType.Kind() {
	case reflect.Interface:
		assertion = assertAssignableToInterface
	default:
		assertion = assertSameType
	}
	return assertion(wantArgType, gotArgType)
}

func assertAssignableToInterface(wantArgType reflect.Type, gotArgType reflect.Type) error {
	if !gotArgType.AssignableTo(wantArgType) {
		return fmt.Errorf("%v is not asignable to %v: %w",
			gotArgType, wantArgType, ErrInvalidArgumentType)
	}
	return nil
}

func assertSameType(wantArgType reflect.Type, gotArgType reflect.Type) error {
	wantArgKind := wantArgType.Kind()
	gotArgKind := gotArgType.Kind()
	if wantArgKind != gotArgKind {
		return fmt.Errorf("got %v, but expected %v: %w",
			gotArgKind, wantArgKind, ErrInvalidArgumentType)
	}
	return nil
}
