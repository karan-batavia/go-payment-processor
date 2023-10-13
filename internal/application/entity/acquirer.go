package entity

import (
	"errors"

	app_error "github.com/sesaquecruz/go-payment-processor/internal/application/errors"
)

var (
	ErrorAcquirerNameIsRequired = errors.New("acquirer name is required")
)

type Acquirer struct {
	Name string
}

func NewAcquirer(name string) *Acquirer {
	return &Acquirer{
		Name: name,
	}
}

func (a *Acquirer) Validate() error {
	errs := make([]error, 0)

	if a.Name == "" {
		errs = append(errs, ErrorAcquirerNameIsRequired)
	}

	if len(errs) > 0 {
		return app_error.NewValidation(errs...)
	}

	return nil
}
