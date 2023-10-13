package entity

import (
	"errors"

	app_error "github.com/sesaquecruz/go-payment-processor/internal/application/errors"
)

var (
	ErrorCardTokenIsRequired      = errors.New("card token is required")
	ErrorCardHolderIsRequired     = errors.New("card holder is required")
	ErrorCardExpirationIsRequired = errors.New("card expiration is required")
	ErrorCardBrandIsRequired      = errors.New("card brand is required")
)

type Card struct {
	Token      string
	Holder     string
	Expiration string
	Brand      string
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
		return app_error.NewValidation(errs...)
	}

	return nil
}
