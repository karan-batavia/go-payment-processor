package entity

import (
	"errors"
	"testing"

	app_error "github.com/sesaquecruz/go-payment-processor/internal/application/errors"
)

func TestCardFactory(t *testing.T) {
	card := NewCard("Token", "Holder", "Expiration", "Brand")

	if card == nil {
		t.Error("card should have been created")
		return
	}

	if card.Token != "Token" {
		t.Error("card token should be Token")
		return
	}

	if card.Holder != "Holder" {
		t.Error("card holder should be Holder")
		return
	}

	if card.Expiration != "Expiration" {
		t.Error("card expiration should be Expiration")
		return
	}

	if card.Brand != "Brand" {
		t.Error("card brand should be Brand")
		return
	}
}

func TestCardValidator(t *testing.T) {
	testCases := []struct {
		Test       string
		Token      string
		Holder     string
		Expiration string
		Brand      string
		errs       []error
	}{
		{
			"token is empty",
			"",
			"Holder",
			"Expiration",
			"Brand",
			[]error{ErrorCardTokenIsRequired},
		},
		{
			"holder is empty",
			"Token",
			"",
			"Expiration",
			"Brand",
			[]error{ErrorCardHolderIsRequired},
		},
		{
			"expiration is empty",
			"Token",
			"Holder",
			"",
			"Brand",
			[]error{ErrorCardExpirationIsRequired},
		},
		{
			"brand is empty",
			"Token",
			"Holder",
			"Expiration",
			"",
			[]error{ErrorCardBrandIsRequired},
		},
		{
			"all fields are invalid",
			"",
			"",
			"",
			"",
			[]error{
				ErrorCardTokenIsRequired,
				ErrorCardHolderIsRequired,
				ErrorCardExpirationIsRequired,
				ErrorCardBrandIsRequired,
			},
		},
		{
			"all fields are valid",
			"Token",
			"Holder",
			"Expiration",
			"Brand",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			err := NewCard(tc.Token, tc.Holder, tc.Expiration, tc.Brand).Validate()

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
