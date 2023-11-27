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
		TestName            string
		StoreIdentification string
		StoreAddress        string
		StoreCep            string
		Err                 *errors.ValidationError
	}{
		{
			"identification is empty",
			"",
			"Address",
			"Cep",
			errors.NewValidationError("store identification is required"),
		},
		{
			"address is empty",
			"Identification",
			"",
			"Cep",
			errors.NewValidationError("store address is required"),
		},
		{
			"cep is empty",
			"Identification",
			"Address",
			"",
			errors.NewValidationError("store cep is required"),
		},
		{
			"all fields are invalid",
			"",
			"",
			"",
			errors.NewValidationError(
				"store identification is required",
				"store address is required",
				"store cep is required",
			),
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
		t.Run(tc.TestName, func(t *testing.T) {
			err := NewStore(tc.StoreIdentification, tc.StoreAddress, tc.StoreCep).Validate()
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
