package entity

import (
	"errors"

	app_errors "github.com/sesaquecruz/go-payment-processor/internal/application/errors"
)

var (
	ErrorPaymentStatusIsRequired = errors.New("payment status is required")
)

type Payment struct {
	Status string
}

func NewPayment(status string) *Payment {
	return &Payment{
		Status: status,
	}
}

func (p *Payment) Validate() error {
	errs := make([]error, 0)

	if p.Status == "" {
		errs = append(errs, ErrorPaymentStatusIsRequired)
	}

	if len(errs) > 0 {
		return app_errors.NewValidation(errs...)
	}

	return nil
}
