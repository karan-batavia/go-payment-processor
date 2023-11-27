package entity

import (
	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
)

type Payment struct {
	Id string
}

func NewPayment(id string) *Payment {
	return &Payment{
		Id: id,
	}
}

func (p *Payment) Validate() error {
	msgs := make([]string, 0)

	if p.Id == "" {
		msgs = append(msgs, "payment id is required")
	}

	if len(msgs) > 0 {
		return errors.NewValidationError(msgs...)
	}

	return nil
}
