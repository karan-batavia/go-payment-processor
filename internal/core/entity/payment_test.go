package entity

import (
	"testing"

	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestPaymentFactory(t *testing.T) {
	payment := NewPayment("Id", "Status")
	assert.NotNil(t, payment)
	assert.Equal(t, payment.Id, "Id")
	assert.Equal(t, payment.Status, "Status")
}

func TestPaymentValidator(t *testing.T) {
	testCases := []struct {
		Test   string
		Id     string
		Status string
		errs   []error
	}{
		{
			"id is empty",
			"",
			"Status",
			[]error{ErrorPaymentIdIsRequired},
		},
		{
			"status is empty",
			"Id",
			"",
			[]error{ErrorPaymentStatusIsRequired},
		},
		{
			"all fields are valid",
			"Id",
			"Status",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			err := NewPayment(tc.Id, tc.Status).Validate()
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
