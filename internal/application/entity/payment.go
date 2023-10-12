package entity

import "github.com/sesaquecruz/go-payment-processor/internal/application/errors"

const (
	PaymentStatusIsRequiredErr = errors.Validation("payment status is required")
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
	if p.Status == "" {
		return PaymentStatusIsRequiredErr
	}

	return nil
}
