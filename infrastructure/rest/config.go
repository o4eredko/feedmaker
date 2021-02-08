package rest

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Config struct {
	Host string
	Port string
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Host, validation.Required, is.Host),
		validation.Field(&c.Port, validation.Required, is.Port),
	)
}
