package entity

import (
	"encoding/json"
	"testing"

	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestTransactionFactory(t *testing.T) {
	card := NewCard("Token", "Holder", "Expiration", "Brand")
	purchase := NewPurchase(9.99, []string{"Item 1", "Item 2"}, 3)
	store := NewStore("Identification", "Address", "Cep")
	acquirer := NewAcquirer("Acquirer")

	transaction := NewTransaction(*card, *purchase, *store, *acquirer)
	assert.NotNil(t, transaction)
	assert.Equal(t, &transaction.Card, card)
	assert.Equal(t, &transaction.Purchase, purchase)
	assert.Equal(t, &transaction.Store, store)
	assert.Equal(t, &transaction.Store, store)
}

func TestTransactionValidator(t *testing.T) {
	testCases := []struct {
		Test     string
		Card     *Card
		Purchase *Purchase
		Store    *Store
		Acquirer *Acquirer
		errs     []error
	}{
		{
			"card is invalid",
			NewCard("Token", "Holder", "", ""),
			NewPurchase(6.99, []string{"Item 1"}, 1),
			NewStore("Identification", "Address", "Cep"),
			NewAcquirer("Acquirer"),
			[]error{
				ErrorCardExpirationIsRequired,
				ErrorCardBrandIsRequired,
			},
		},
		{
			"purchase is invalid",
			NewCard("Token", "Holder", "Expiration", "Brand"),
			NewPurchase(0, []string{"Item 1"}, -1),
			NewStore("Identification", "Address", "Cep"),
			NewAcquirer("Acquirer"),
			[]error{
				ErrorPurchaseValueIsInvalid,
				ErrorPurchaseInstallmentsIsInvalid,
			},
		},
		{
			"store is invalid",
			NewCard("Token", "Holder", "Expiration", "Brand"),
			NewPurchase(6.99, []string{"Item 1"}, 1),
			NewStore("", "", ""),
			NewAcquirer("Acquirer"),
			[]error{
				ErrorStoreIdentificationIsRequired,
				ErrorStoreAddressIsRequired,
				ErrorStoreCepIsRequired,
			},
		},
		{
			"acquirer is invalid",
			NewCard("Token", "Holder", "Expiration", "Brand"),
			NewPurchase(6.99, []string{"Item 1"}, 1),
			NewStore("Identification", "Address", "Cep"),
			NewAcquirer(""),
			[]error{
				ErrorAcquirerNameIsRequired,
			},
		},
		{
			"all fields are invalid",
			NewCard("Token", "Holder", "", "Brand"),
			NewPurchase(6.99, []string{}, 1),
			NewStore("", "Address", "Cep"),
			NewAcquirer(""),
			[]error{
				ErrorCardExpirationIsRequired,
				ErrorPurchaseItemsIsRequired,
				ErrorStoreIdentificationIsRequired,
				ErrorAcquirerNameIsRequired,
			},
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
		t.Run(tc.Test, func(t *testing.T) {
			err := NewTransaction(*tc.Card, *tc.Purchase, *tc.Store, *tc.Acquirer).Validate()
			if tc.errs == nil && err == nil {
				return
			}

			var v *errors.ValidationError
			assert.ErrorAs(t, err, &v)

			errs := v.Unwrap()
			assert.Equal(t, len(tc.errs), len(errs))

			for i, err := range tc.errs {
				assert.ErrorIs(t, err, errs[i])
			}
		})
	}
}

func TestMarshalTransaction(t *testing.T) {
	expected := `
		{
			"card": {
				"token": "Token",
				"holder": "Holder",
				"expiration": "Expiration",
				"brand": "Brand"
			},
			"purchase": {
				"value": 4.99,
				"items": ["Item 1, Item 2"],
				"installments": 3
			},
			"store": {
				"identification": "Identification",
				"address": "Address",
				"cep": "Cep"
			}
		}
	`

	transaction := Transaction{
		Card: Card{
			Token:      "Token",
			Holder:     "Holder",
			Expiration: "Expiration",
			Brand:      "Brand",
		},
		Purchase: Purchase{
			Value:        4.99,
			Items:        []string{"Item 1, Item 2"},
			Installments: 3,
		},
		Store: Store{
			Identification: "Identification",
			Address:        "Address",
			Cep:            "Cep",
		},
		Acquirer: Acquirer{
			Name: "Name",
		},
	}

	actual, err := json.Marshal(&transaction)
	assert.Nil(t, err)
	assert.JSONEq(t, expected, string(actual))
}
