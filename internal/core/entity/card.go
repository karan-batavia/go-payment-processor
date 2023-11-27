package entity

import (
	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
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
	msgs := make([]string, 0)

	if c.Token == "" {
		msgs = append(msgs, "card token is required")
	}

	if c.Holder == "" {
		msgs = append(msgs, "card holder is required")
	}

	if c.Expiration == "" {
		msgs = append(msgs, "card expiration is required")
	}

	if c.Brand == "" {
		msgs = append(msgs, "card brand is required")
	}

	if len(msgs) > 0 {
		return errors.NewValidationError(msgs...)
	}

	return nil
}
