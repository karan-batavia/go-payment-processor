package entity

import (
	app_errors "github.com/sesaquecruz/go-payment-processor/internal/application/errors"
)

const (
	ErrorPaymentIdIsRequired     = app_errors.Error("payment id is required")
	ErrorPaymentStatusIsRequired = app_errors.Error("payment status is required")
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
		return app_errors.NewValidation(errs...)
	}

	return nil
}
