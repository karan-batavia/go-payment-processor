package entity

import (
	"testing"

	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestPaymentFactory(t *testing.T) {
	payment := NewPayment("Id")
	assert.NotNil(t, payment)
	assert.Equal(t, payment.Id, "Id")
}

func TestPaymentValidator(t *testing.T) {
	testCases := []struct {
		TestName  string
		PaymentId string
		Err       *errors.ValidationError
	}{
		{
			"id is empty",
			"",
			errors.NewValidationError("payment id is required"),
		},
		{
			"all fields are valid",
			"Id",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.TestName, func(t *testing.T) {
			err := NewPayment(tc.PaymentId).Validate()
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
