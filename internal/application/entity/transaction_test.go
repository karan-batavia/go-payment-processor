package entity

import (
	"errors"
	"reflect"
	"testing"

	app_error "github.com/sesaquecruz/go-payment-processor/internal/application/errors"
)

func TestTransactionFactory(t *testing.T) {
	card := NewCard("Token", "Holder", "Expiration", "Brand")
	purchase := NewPurchase(9.99, []string{"Item 1", "Item 2"}, 3)
	store := NewStore("Identification", "Address", "Cep")
	acquirer := NewAcquirer("Acquirer")
	transaction := NewTransaction(*card, *purchase, *store, *acquirer)

	if transaction == nil {
		t.Error("transaction should have been created")
		return
	}

	if !reflect.DeepEqual(&transaction.Card, card) {
		t.Errorf("expected: %v, got: %v", transaction.Card, card)
		return
	}

	if !reflect.DeepEqual(&transaction.Purchase, purchase) {
		t.Errorf("expected: %v, got: %v", transaction.Purchase, purchase)
		return
	}

	if !reflect.DeepEqual(&transaction.Store, store) {
		t.Errorf("expected: %v, got: %v", transaction.Store, store)
		return
	}

	if !reflect.DeepEqual(&transaction.Acquirer, acquirer) {
		t.Errorf("expected: %v, got: %v", transaction.Acquirer, acquirer)
		return
	}
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

			if tc.errs == nil && err != nil {
				t.Errorf("expected: %v, got: %v", tc.errs, err)
				return
			}

			var v *app_error.Validation
			if !errors.As(err, &v) {
				t.Errorf("expected a validation error, got: %v", err)
				return
			}

			errs := v.Unwrap()

			if len(tc.errs) != len(errs) {
				t.Errorf("expected %d errors, got: %d errors", len(tc.errs), len(errs))
				return
			}

			for i, err := range tc.errs {
				if !errors.Is(err, errs[i]) {
					t.Errorf("expected: %v, got: %v", err, errs[i])
					return
				}
			}
		})
	}
}
