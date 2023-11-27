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
		TestName       string
		CardToken      string
		CardHolder     string
		CardExpiration string
		CardBrand      string
		Err            *errors.ValidationError
	}{
		{
			"token is empty",
			"",
			"Holder",
			"Expiration",
			"Brand",
			errors.NewValidationError("card token is required"),
		},
		{
			"holder is empty",
			"Token",
			"",
			"Expiration",
			"Brand",
			errors.NewValidationError("card holder is required"),
		},
		{
			"expiration is empty",
			"Token",
			"Holder",
			"",
			"Brand",
			errors.NewValidationError("card expiration is required"),
		},
		{
			"brand is empty",
			"Token",
			"Holder",
			"Expiration",
			"",
			errors.NewValidationError("card brand is required"),
		},
		{
			"all fields are invalid",
			"",
			"",
			"",
			"",
			errors.NewValidationError(
				"card token is required",
				"card holder is required",
				"card expiration is required",
				"card brand is required",
			),
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
		t.Run(tc.TestName, func(t *testing.T) {
			err := NewCard(tc.CardToken, tc.CardHolder, tc.CardExpiration, tc.CardBrand).Validate()
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
