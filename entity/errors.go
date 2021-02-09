package entity

import "errors"

var (
	ErrInvalidGenerationType = errors.New("this generation type is invalid")
	ErrInvalidTimestamp      = errors.New("invalid timestamp")
)
