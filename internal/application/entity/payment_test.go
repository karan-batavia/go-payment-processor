package entity

import (
	"errors"
	"testing"

	app_error "github.com/sesaquecruz/go-payment-processor/internal/application/errors"
)

func TestPaymentFactory(t *testing.T) {
	payment := NewPayment("Status")

	if payment == nil {
		t.Error("payment should have been created")
		return
	}

	if payment.Status != "Status" {
		t.Error("payment status should be Status")
		return
	}
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
