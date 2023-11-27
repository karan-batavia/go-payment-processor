package entity

import (
	"testing"

	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestTransactionFactory(t *testing.T) {
	card := NewCard("Token", "Holder", "Expiration", "Brand")
	purchase := NewPurchase(9.99, []string{"Item 1", "Item 2"}, 3)
	store := NewStore("Identification", "Address", "Cep")
	acquirer := NewAcquirer("Acquirer")

	transaction := NewTransaction(card, purchase, store, acquirer)
	assert.NotNil(t, transaction)
	assert.Equal(t, transaction.Card, card)
	assert.Equal(t, transaction.Purchase, purchase)
	assert.Equal(t, transaction.Store, store)
	assert.Equal(t, transaction.Store, store)
}

func TestTransactionValidator(t *testing.T) {
	testCases := []struct {
		TestName string
		Card     *Card
		Purchase *Purchase
		Store    *Store
		Acquirer *Acquirer
		Err      *errors.ValidationError
	}{
		{
			"card is invalid",
			NewCard("Token", "Holder", "", ""),
			NewPurchase(6.99, []string{"Item 1"}, 1),
			NewStore("Identification", "Address", "Cep"),
			NewAcquirer("Acquirer"),
			errors.NewValidationError(
				"card expiration is required",
				"card brand is required",
			),
		},
		{
			"purchase is invalid",
			NewCard("Token", "Holder", "Expiration", "Brand"),
			NewPurchase(0, []string{"Item 1"}, -1),
			NewStore("Identification", "Address", "Cep"),
			NewAcquirer("Acquirer"),
			errors.NewValidationError(
				"purchase value is invalid",
				"purchase installments is invalid",
			),
		},
		{
			"store is invalid",
			NewCard("Token", "Holder", "Expiration", "Brand"),
			NewPurchase(6.99, []string{"Item 1"}, 1),
			NewStore("", "", ""),
			NewAcquirer("Acquirer"),
			errors.NewValidationError(
				"store identification is required",
				"store address is required",
				"store cep is required",
			),
		},
		{
			"acquirer is invalid",
			NewCard("Token", "Holder", "Expiration", "Brand"),
			NewPurchase(6.99, []string{"Item 1"}, 1),
			NewStore("Identification", "Address", "Cep"),
			NewAcquirer(""),
			errors.NewValidationError("acquirer name is required"),
		},
		{
			"all fields are invalid",
			NewCard("Token", "Holder", "", "Brand"),
			NewPurchase(6.99, []string{}, 1),
			NewStore("", "Address", "Cep"),
			NewAcquirer(""),
			errors.NewValidationError(
				"card expiration is required",
				"purchase items is required",
				"store identification is required",
				"acquirer name is required",
			),
		},
		{
			"all fields are valid",
			NewCard("Token", "Holder", "Expiration", "Brand"),
			NewPurchase(6.99, []string{"Item 1"}, 1),
			NewStore("Identification", "Address", "Cep"),
			NewAcquirer("Acquirer"),
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.TestName, func(t *testing.T) {
			err := NewTransaction(tc.Card, tc.Purchase, tc.Store, tc.Acquirer).Validate()
			if tc.Err == nil && err == nil {
				return
			}

			var verr *errors.ValidationError
			assert.ErrorAs(t, err, &verr)
			assert.Equal(t, len(tc.Err.Messages), len(verr.Messages))

			for i, msg := range tc.Err.Messages {
				assert.Equal(t, msg, verr.Messages[i])
			}
		})
	}
}
