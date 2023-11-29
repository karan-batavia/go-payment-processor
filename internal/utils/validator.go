package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func GetValidator() *validator.Validate {
	return validate
}
