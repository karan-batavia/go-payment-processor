package entity

import (
	"errors"

	core_errors "github.com/sesaquecruz/go-payment-processor/internal/core/errors"
)

type Transaction struct {
	Card     *Card     `json:"card"`
	Purchase *Purchase `json:"purchase"`
	Store    *Store    `json:"store"`
	Acquirer *Acquirer `json:"-"`
}

func NewTransaction(card *Card, purchase *Purchase, store *Store, acquirer *Acquirer) *Transaction {
	return &Transaction{
		Card:     card,
		Purchase: purchase,
		Store:    store,
		Acquirer: acquirer,
	}
}

func (t *Transaction) Validate() error {
	msgs := make([]string, 0)

	err := t.Card.Validate()
	if err != nil {
		var v *core_errors.ValidationError
		if errors.As(err, &v) {
			msgs = append(msgs, v.Messages...)
		}
	}

	err = t.Purchase.Validate()
	if err != nil {
		var v *core_errors.ValidationError
		if errors.As(err, &v) {
			msgs = append(msgs, v.Messages...)
		}
	}

	err = t.Store.Validate()
	if err != nil {
		var v *core_errors.ValidationError
		if errors.As(err, &v) {
			msgs = append(msgs, v.Messages...)
		}
	}

	err = t.Acquirer.Validate()
	if err != nil {
		var v *core_errors.ValidationError
		if errors.As(err, &v) {
			msgs = append(msgs, v.Messages...)
		}
	}

	if len(msgs) > 0 {
		return core_errors.NewValidationError(msgs...)
	}

	return nil
}
