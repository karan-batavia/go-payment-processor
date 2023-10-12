package entity

import (
	"errors"
	"reflect"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	card := NewCard("Token", "Holder", "Expiration", "Brand")
	purchase := NewPurchase(9.99, []string{"Item 1", "Item 2"}, 3)
	store := NewStore("Identification", "Address", "Cep")
	transaction := NewTransaction(card, purchase, store, "Acquirer")

	if transaction == nil {
		t.Error("transaction should have been created")
		return
	}

	if !reflect.DeepEqual(transaction.Card, card) {
		t.Errorf("transaction card should be %v", card)
	}

	if !reflect.DeepEqual(transaction.Purchase, purchase) {
		t.Errorf("transaction purchase should be %v", purchase)
	}

	if !reflect.DeepEqual(transaction.Store, store) {
		t.Errorf("transaction store should be %v", store)
	}

	if transaction.Acquirer != "Acquirer" {
		t.Error("transaction acquirer should be Acquirer")
	}
}

func TestTransactionFields(t *testing.T) {
	testCases := []struct {
		Test     string
		Card     *Card
		Purchase *Purchase
		Store    *Store
		Acquirer string
		err      error
	}{
		{
			"card is nil",
			nil,
			NewPurchase(9.99, []string{"Item 1", "Item 2"}, 3),
			NewStore("Identification", "Address", "Cep"),
			"Acquirer",
			TransactionCardIsRequired,
		},
		{
			"card is invalid",
			NewCard("Token", "Holder", "Expiratino", ""),
			NewPurchase(9.99, []string{"Item 1", "Item 2"}, 3),
			NewStore("Identification", "Address", "Cep"),
			"Acquirer",
			CardBrandIsRequiredErr,
		},
		{
			"purchase is nil",
			NewCard("Token", "Holder", "Expiratino", "Brand"),
			nil,
			NewStore("Identification", "Address", "Cep"),
			"Acquirer",
			TransactionPurchaseIsRequired,
		},
		{
			"purchase is invalid",
			NewCard("Token", "Holder", "Expiratino", "Brand"),
			NewPurchase(9.99, []string{"Item 1", "Item 2"}, 0),
			NewStore("Identification", "Address", "Cep"),
			"Acquirer",
			PurchaseInstallmentsIsInvalidErr,
		},
		{
			"store is nil",
			NewCard("Token", "Holder", "Expiratino", "Brand"),
			NewPurchase(9.99, []string{"Item 1", "Item 2"}, 3),
			nil,
			"Acquirer",
			TransactionStoreIsRequired,
		},
		{
			"store is invalid",
			NewCard("Token", "Holder", "Expiratino", "Brand"),
			NewPurchase(9.99, []string{"Item 1", "Item 2"}, 3),
			NewStore("Identification", "Address", ""),
			"Acquirer",
			StoreCepIsRequiredErr,
		},
		{
			"acquirer is empty",
			NewCard("Token", "Holder", "Expiratino", "Brand"),
			NewPurchase(9.99, []string{"Item 1", "Item 2"}, 3),
			NewStore("Identification", "Address", "Cep"),
			"",
			TransactionAcquirerIsRequired,
		},
		{
			"all fields are valid",
			NewCard("Token", "Holder", "Expiratino", "Brand"),
			NewPurchase(9.99, []string{"Item 1", "Item 2"}, 3),
			NewStore("Identification", "Address", "Cep"),
			"Acquirer",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			err := NewTransaction(tc.Card, tc.Purchase, tc.Store, tc.Acquirer).Validate()
			if !errors.Is(err, tc.err) {
				t.Errorf("expected: %v, got: %v", tc.err, err)
			}
		})
	}
}
