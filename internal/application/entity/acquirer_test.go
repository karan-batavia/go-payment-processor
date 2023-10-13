package entity

import (
	"errors"
	"testing"

	app_error "github.com/sesaquecruz/go-payment-processor/internal/application/errors"
)

func TestAcquirerFactory(t *testing.T) {
	acquirer := NewAcquirer("Acquirer")

	if acquirer == nil {
		t.Error("acquirer should have been created")
		return
	}

	if acquirer.Name != "Acquirer" {
		t.Error("acquirer name should be Acquirer")
		return
	}
}

func TestAcquirerValidator(t *testing.T) {
	testCases := []struct {
		Test string
		Name string
		errs []error
	}{
		{
			"name is emtpy",
			"",
			[]error{ErrorAcquirerNameIsRequired},
		},
		{
			"all fields are valid",
			"Acquirer",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := NewAcquirer(tc.Name).Validate()

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
