package entity

import (
	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
)

type Purchase struct {
	Value        float64
	Items        []string
	Installments int
}

func NewPurchase(value float64, items []string, installments int) *Purchase {
	return &Purchase{
		Value:        value,
		Items:        items,
		Installments: installments,
	}
}

func (p *Purchase) Validate() error {
	msgs := make([]string, 0)

	if p.Value <= 0 {
		msgs = append(msgs, "purchase value is invalid")
	}

	if p.Items == nil || len(p.Items) == 0 {
		msgs = append(msgs, "purchase items is required")
	} else {
		for _, item := range p.Items {
			if item == "" {
				msgs = append(msgs, "purchase items is invalid")
				break
			}
		}
	}

	if p.Installments <= 0 {
		msgs = append(msgs, "purchase installments is invalid")
	}

	if len(msgs) > 0 {
		return errors.NewValidationError(msgs...)
	}

	return nil
}
