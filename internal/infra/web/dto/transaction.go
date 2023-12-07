package dto

import (
	"fmt"

	web_errors "github.com/sesaquecruz/go-payment-processor/internal/infra/web/errors"
	"github.com/sesaquecruz/go-payment-processor/internal/utils"

	"github.com/go-playground/validator/v10"
)

type Transaction struct {
	CardToken            string   `json:"card_token"            validate:"required"`
	PurchaseValue        float64  `json:"purchase_value"        validate:"required"`
	PurchaseItens        []string `json:"purchase_items"        validate:"required"`
	PurchaseInstallments int      `json:"purchase_installments" validate:"required"`
	StoreIdentification  string   `json:"store_identification"  validate:"required"`
	StoreAddress         string   `json:"store_address"         validate:"required"`
	StoreCep             string   `json:"store_cep"             validate:"required"`
	AcquirerName         string   `json:"acquirer_name"         validate:"required"`
}

func (t *Transaction) Validate() error {
	err := utils.GetValidator().Struct(t)
	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)
	msgs := make([]string, 0, len(errs))

	for _, e := range errs {
		msg := fmt.Sprintf("%s is required", utils.GetNamespaceError(e))
		msgs = append(msgs, msg)
	}

	return web_errors.NewError(msgs...)
}
