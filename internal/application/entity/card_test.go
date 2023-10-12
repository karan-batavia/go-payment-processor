package entity

import (
	"errors"
	"testing"
)

func TestCreateCard(t *testing.T) {
	card := NewCard("Token", "Holder", "Expiration", "Brand")

	if card == nil {
		t.Error("card should have been created")
		return
	}

	if card.Token != "Token" {
		t.Error("card token should be Token")
	}

	if card.Holder != "Holder" {
		t.Error("card holder should be Holder")
	}

	if card.Expiration != "Expiration" {
		t.Error("card expiration should be Expiration")
	}

	if card.Brand != "Brand" {
		t.Error("card brand should be Brand")
	}
}

func TestCardFields(t *testing.T) {
	testCases := []struct {
		Test       string
		Token      string
		Holder     string
		Expiration string
		Brand      string
		err        error
	}{
		{
			"token is empty",
			"",
			"Holder",
			"Expiration",
			"Brand",
			CardTokenIsRequiredErr,
		},
		{
			"holder is empty",
			"Token",
			"",
			"Expiration",
			"Brand",
			CardHolderIsRequiredErr,
		},
		{
			"expiration is empty",
			"Token",
			"Holder",
			"",
			"Brand",
			CardExpirationIsRequiredErr,
		},
		{
			"brand is empty",
			"Token",
			"Holder",
			"Expiration",
			"",
			CardBrandIsRequiredErr,
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
			if !errors.Is(err, tc.err) {
				t.Errorf("expected: %v, got: %v", tc.err, err)
			}
		})
	}
}
