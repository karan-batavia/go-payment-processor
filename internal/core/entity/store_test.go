package entity

import (
	"testing"

	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreateStore(t *testing.T) {
	store := NewStore("Identification", "Address", "Cep")
	assert.NotNil(t, store)
	assert.Equal(t, store.Identification, "Identification")
	assert.Equal(t, store.Address, "Address")
	assert.Equal(t, store.Cep, "Cep")
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
