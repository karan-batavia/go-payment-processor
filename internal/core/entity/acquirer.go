package entity

import (
	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
)

const (
	ErrorAcquirerNameIsRequired = errors.Error("acquirer name is required")
)

type Acquirer struct {
	Name string `json:"name"`
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
		return errors.NewValidationError(errs...)
	}

	return nil
}
