package entity

import (
	"testing"

	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestPurchaseFactory(t *testing.T) {
	items := []string{"Item 1", "Item 2"}

	purchase := NewPurchase(9.99, items, 5)
	assert.NotNil(t, purchase)
	assert.Equal(t, purchase.Value, 9.99)
	assert.EqualValues(t, purchase.Items, items)
	assert.Equal(t, purchase.Installments, 5)
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
