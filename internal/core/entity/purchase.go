package entity

import (
	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
)

const (
	ErrorPurchaseItemsIsRequired = errors.Error("purchase items is required")

	ErrorPurchaseValueIsInvalid        = errors.Error("purchase value is invalid")
	ErrorPurchaseItemsIsInvalid        = errors.Error("purchase items is invalid")
	ErrorPurchaseInstallmentsIsInvalid = errors.Error("purchase installments is invalid")
)

type Purchase struct {
	Value        float64  `json:"value"`
	Items        []string `json:"items"`
	Installments int      `json:"installments"`
}

func NewPurchase(value float64, items []string, installments int) *Purchase {
	return &Purchase{
		Value:        value,
		Items:        items,
		Installments: installments,
	}
}

func (p *Purchase) Validate() error {
	errs := make([]error, 0)

	if p.Value <= 0 {
		errs = append(errs, ErrorPurchaseValueIsInvalid)
	}

	if p.Items == nil || len(p.Items) == 0 {
		errs = append(errs, ErrorPurchaseItemsIsRequired)
	} else {
		for _, item := range p.Items {
			if item == "" {
				errs = append(errs, ErrorPurchaseItemsIsInvalid)
				break
			}
		}
	}

	if p.Installments <= 0 {
		errs = append(errs, ErrorPurchaseInstallmentsIsInvalid)
	}

	if len(errs) > 0 {
		return errors.NewValidationError(errs...)
	}

	return nil
}
