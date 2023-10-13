package entity

import (
	"errors"
	"reflect"
	"testing"

	app_error "github.com/sesaquecruz/go-payment-processor/internal/application/errors"
)

func TestPurchaseFactory(t *testing.T) {
	items := []string{"Item 1", "Item 2"}
	purchase := NewPurchase(9.99, items, 5)

	if purchase == nil {
		t.Error("purchase should have been created")
		return
	}

	if purchase.Value != 9.99 {
		t.Error("purchase value should be 9.99")
		return
	}

	if !reflect.DeepEqual(purchase.Items, items) {
		t.Errorf("purchase items should be %v", items)
		return
	}

	if purchase.Installments != 5 {
		t.Error("purchase installments shoudl be 5")
		return
	}
}

func TestPurchaseValidator(t *testing.T) {
	testCases := []struct {
		Test         string
		Value        float64
		Items        []string
		Installments int
		errs         []error
	}{
		{
			"value is negative",
			-1,
			[]string{"Item 1", "Item 2"},
			1,
			[]error{ErrorPurchaseValueIsInvalid},
		},
		{
			"value is zero",
			0,
			[]string{"Item 1", "Item 2"},
			1,
			[]error{ErrorPurchaseValueIsInvalid},
		},
		{
			"items is nil",
			1,
			nil,
			1,
			[]error{ErrorPurchaseItemsIsRequired},
		},
		{
			"items is empty",
			1,
			[]string{},
			1,
			[]error{ErrorPurchaseItemsIsRequired},
		},
		{
			"items has empty elements",
			1,
			[]string{"Item 1", ""},
			1,
			[]error{ErrorPurchaseItemsIsInvalid},
		},
		{
			"installments is negative",
			1,
			[]string{"Item 1", "Item 2"},
			-1,
			[]error{ErrorPurchaseInstallmentsIsInvalid},
		},
		{
			"installments is zero",
			1,
			[]string{"Item 1", "Item 2"},
			0,
			[]error{ErrorPurchaseInstallmentsIsInvalid},
		},
		{
			"all fields are invalid",
			0,
			[]string{""},
			0,
			[]error{
				ErrorPurchaseValueIsInvalid,
				ErrorPurchaseItemsIsInvalid,
				ErrorPurchaseInstallmentsIsInvalid,
			},
		},
		{
			"all fields are valid",
			1,
			[]string{"Item 1", "Item 2"},
			1,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			err := NewPurchase(tc.Value, tc.Items, tc.Installments).Validate()

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
