package entity

import (
	"errors"
	"testing"
)

func TestCreateStore(t *testing.T) {
	store := NewStore("Identification", "Address", "Cep")

	if store == nil {
		t.Error("store should have been created")
		return
	}

	if store.Identification != "Identification" {
		t.Error("store identification should be Identification")
	}

	if store.Address != "Address" {
		t.Error("store address should be Address")
	}

	if store.Cep != "Cep" {
		t.Error("store cep should be Cep")
	}
}

func TestStoreFields(t *testing.T) {
	testCase := []struct {
		Test           string
		Identification string
		Address        string
		Cep            string
		err            error
	}{
		{
			"identification is empty",
			"",
			"Address",
			"Cep",
			StoreIdentificationIsRequiredErr,
		},
		{
			"address is empty",
			"Identification",
			"",
			"",
			StoreAddressIsRequiredErr,
		},
		{
			"cep is empty",
			"Identification",
			"Address",
			"",
			StoreCepIsRequiredErr,
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
			if !errors.Is(err, tc.err) {
				t.Errorf("expected: %v, got: %v", tc.err, err)
			}
		})
	}
}
