package entity

import (
	"errors"
	"testing"

	app_error "github.com/sesaquecruz/go-payment-processor/internal/application/errors"
)

func TestCreateStore(t *testing.T) {
	store := NewStore("Identification", "Address", "Cep")

	if store == nil {
		t.Error("store should have been created")
		return
	}

	if store.Identification != "Identification" {
		t.Error("store identification should be Identification")
		return
	}

	if store.Address != "Address" {
		t.Error("store address should be Address")
		return
	}

	if store.Cep != "Cep" {
		t.Error("store cep should be Cep")
		return
	}
}

func TestStoreValidator(t *testing.T) {
	testCase := []struct {
		Test           string
		Identification string
		Address        string
		Cep            string
		errs           []error
	}{
		{
			"identification is empty",
			"",
			"Address",
			"Cep",
			[]error{ErrorStoreIdentificationIsRequired},
		},
		{
			"address is empty",
			"Identification",
			"",
			"Cep",
			[]error{ErrorStoreAddressIsRequired},
		},
		{
			"cep is empty",
			"Identification",
			"Address",
			"",
			[]error{ErrorStoreCepIsRequired},
		},
		{
			"all fields are invalid",
			"",
			"",
			"",
			[]error{
				ErrorStoreIdentificationIsRequired,
				ErrorStoreAddressIsRequired,
				ErrorStoreCepIsRequired,
			},
		},
		{
			"all fields are valid",
			"Identification",
			"Address",
			"Cep",
			nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.Test, func(t *testing.T) {
			err := NewStore(tc.Identification, tc.Address, tc.Cep).Validate()

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
