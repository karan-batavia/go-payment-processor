package entity

import (
	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
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
	msgs := make([]string, 0)

	if a.Name == "" {
		msgs = append(msgs, "acquirer name is required")
	}

	if len(msgs) > 0 {
		return errors.NewValidationError(msgs...)
	}

	return nil
}
