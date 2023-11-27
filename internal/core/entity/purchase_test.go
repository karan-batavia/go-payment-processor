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
		TestName             string
		PurchaseValue        float64
		PurchaseItems        []string
		PurchaseInstallments int
		Err                  *errors.ValidationError
	}{
		{
			"value is negative",
			-1,
			[]string{"Item 1", "Item 2"},
			1,
			errors.NewValidationError("purchase value is invalid"),
		},
		{
			"value is zero",
			0,
			[]string{"Item 1", "Item 2"},
			1,
			errors.NewValidationError("purchase value is invalid"),
		},
		{
			"items is nil",
			1,
			nil,
			1,
			errors.NewValidationError("purchase items is required"),
		},
		{
			"items is empty",
			1,
			[]string{},
			1,
			errors.NewValidationError("purchase items is required"),
		},
		{
			"items has empty elements",
			1,
			[]string{"Item 1", ""},
			1,
			errors.NewValidationError("purchase items is invalid"),
		},
		{
			"installments is negative",
			1,
			[]string{"Item 1", "Item 2"},
			-1,
			errors.NewValidationError("purchase installments is invalid"),
		},
		{
			"installments is zero",
			1,
			[]string{"Item 1", "Item 2"},
			0,
			errors.NewValidationError("purchase installments is invalid"),
		},
		{
			"all fields are invalid",
			0,
			[]string{""},
			0,
			errors.NewValidationError(
				"purchase value is invalid",
				"purchase items is invalid",
				"purchase installments is invalid",
			),
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
		t.Run(tc.TestName, func(t *testing.T) {
			err := NewPurchase(tc.PurchaseValue, tc.PurchaseItems, tc.PurchaseInstallments).Validate()
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
