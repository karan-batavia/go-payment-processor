package entity

import (
	"testing"

	app_errors "github.com/sesaquecruz/go-payment-processor/internal/application/errors"

	"github.com/stretchr/testify/assert"
)

func TestPaymentFactory(t *testing.T) {
	payment := NewPayment("Status")
	assert.NotNil(t, payment)
	assert.Equal(t, payment.Status, "Status")
}

func TestPaymentValidator(t *testing.T) {
	testCases := []struct {
		Test   string
		Status string
		errs   []error
	}{
		{
			"status is empty",
			"",
			[]error{ErrorPaymentStatusIsRequired},
		},
		{
			"all fields are valid",
			"Status",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			err := NewPayment(tc.Status).Validate()
			if tc.errs == nil && err == nil {
				return
			}

			var v *app_errors.Validation
			assert.ErrorAs(t, err, &v)

			errs := v.Unwrap()
			assert.Equal(t, len(tc.errs), len(errs))

			for i, err := range tc.errs {
				assert.ErrorIs(t, err, errs[i])
			}
		})
	}
}
