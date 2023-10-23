package entity

import (
	app_errors "github.com/sesaquecruz/go-payment-processor/internal/application/errors"
)

const (
	ErrorCardTokenIsRequired      = app_errors.Error("card token is required")
	ErrorCardHolderIsRequired     = app_errors.Error("card holder is required")
	ErrorCardExpirationIsRequired = app_errors.Error("card expiration is required")
	ErrorCardBrandIsRequired      = app_errors.Error("card brand is required")
)

type Card struct {
	Token      string `json:"token"`
	Holder     string `json:"holder"`
	Expiration string `json:"expiration"`
	Brand      string `json:"brand"`
}

func NewCard(token string, holder string, expiration string, brand string) *Card {
	return &Card{
		Token:      token,
		Holder:     holder,
		Expiration: expiration,
		Brand:      brand,
	}
}

func (c *Card) Validate() error {
	errs := make([]error, 0)

	if c.Token == "" {
		errs = append(errs, ErrorCardTokenIsRequired)
	}

	if c.Holder == "" {
		errs = append(errs, ErrorCardHolderIsRequired)
	}

	if c.Expiration == "" {
		errs = append(errs, ErrorCardExpirationIsRequired)
	}

	if c.Brand == "" {
		errs = append(errs, ErrorCardBrandIsRequired)
	}

	if len(errs) > 0 {
		return app_errors.NewValidation(errs...)
	}

	return nil
}
