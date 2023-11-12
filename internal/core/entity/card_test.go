package entity

import (
	"testing"

	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestCardFactory(t *testing.T) {
	card := NewCard("Token", "Holder", "Expiration", "Brand")
	assert.NotNil(t, card)
	assert.Equal(t, card.Token, "Token")
	assert.Equal(t, card.Holder, "Holder")
	assert.Equal(t, card.Expiration, "Expiration")
	assert.Equal(t, card.Brand, "Brand")
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
