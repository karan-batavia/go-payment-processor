package entity

import (
	app_errors "github.com/sesaquecruz/go-payment-processor/internal/application/errors"
)

const (
	ErrorAcquirerNameIsRequired = app_errors.Error("acquirer name is required")
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
		return app_errors.NewValidation(errs...)
	}

	return nil
}
