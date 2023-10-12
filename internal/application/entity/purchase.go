package entity

import "github.com/sesaquecruz/go-payment-processor/internal/application/errors"

const (
	PurchaseItemsIsRequiredErr = errors.Validation("purchase items is required")

	PurchaseValueIsInvalidErr        = errors.Validation("purchase value is invalid")
	PurchaseItemsAreInvalidErr       = errors.Validation("purchase items are invalid")
	PurchaseInstallmentsIsInvalidErr = errors.Validation("purchase installments is invalid")
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
	if p.Value <= 0 {
		return PurchaseValueIsInvalidErr
	}

	if p.Items == nil || len(p.Items) == 0 {
		return PurchaseItemsIsRequiredErr
	}

	for _, item := range p.Items {
		if item == "" {
			return PurchaseItemsAreInvalidErr
		}
	}

	if p.Installments <= 0 {
		return PurchaseInstallmentsIsInvalidErr
	}

	return nil
}
