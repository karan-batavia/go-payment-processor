package entity

import (
	"errors"
	"testing"
)

func TestCreatePayment(t *testing.T) {
	payment := NewPayment("Status")

	if payment == nil {
		t.Error("payment should have been created")
		return
	}

	if payment.Status != "Status" {
		t.Error("payment status should be Status")
	}
}

func TestPaymentFields(t *testing.T) {
	testCases := []struct {
		Test   string
		Status string
		err    error
	}{
		{
			"status is empty",
			"",
			PaymentStatusIsRequiredErr,
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
			if !errors.Is(err, tc.err) {
				t.Errorf("expected: %v, got: %v", tc.err, err)
			}
		})
	}
}
