package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func GetValidator() *validator.Validate {
	return validate
}

func GetNamespaceError(err validator.FieldError) string {
	namespace := err.Namespace()
	msg := strings.Builder{}

	for i, c := range namespace {
		if c == '.' {
			continue
		} else if c > 'Z' {
			msg.WriteRune(c)
		} else {
			if i > 0 {
				msg.WriteRune(' ')
			}
			msg.WriteRune(c + 32)
		}
	}

	return msg.String()
}
