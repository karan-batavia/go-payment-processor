package entity

import "github.com/sesaquecruz/go-payment-processor/internal/application/errors"

const (
	CardTokenIsRequiredErr      = errors.Validation("card token is required")
	CardHolderIsRequiredErr     = errors.Validation("card holder is required")
	CardExpirationIsRequiredErr = errors.Validation("card expiration is required")
	CardBrandIsRequiredErr      = errors.Validation("card brand is required")
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
	if c.Token == "" {
		return CardTokenIsRequiredErr
	}

	if c.Holder == "" {
		return CardHolderIsRequiredErr
	}

	if c.Expiration == "" {
		return CardExpirationIsRequiredErr
	}

	if c.Brand == "" {
		return CardBrandIsRequiredErr
	}

	return nil
}
