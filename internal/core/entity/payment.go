package entity

import (
	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
)

const (
	ErrorPaymentIdIsRequired     = errors.Error("payment id is required")
	ErrorPaymentStatusIsRequired = errors.Error("payment status is required")
)

const (
	PaymentStatusPaid   = "paid"
	PaymentStatusFailed = "failed"
)

type Payment struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

func NewPayment(id string, status string) *Payment {
	return &Payment{
		Id:     id,
		Status: status,
	}
}

func (p *Payment) Validate() error {
	errs := make([]error, 0)

	if p.Id == "" {
		errs = append(errs, ErrorPaymentIdIsRequired)
	}

	if p.Status == "" {
		errs = append(errs, ErrorPaymentStatusIsRequired)
	}

	if len(errs) > 0 {
		return errors.NewValidationError(errs...)
	}

	return nil
}
