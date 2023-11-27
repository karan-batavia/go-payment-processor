package entity

import (
	"testing"

	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestAcquirerFactory(t *testing.T) {
	acquirer := NewAcquirer("Acquirer")
	assert.NotNil(t, acquirer)
	assert.Equal(t, acquirer.Name, "Acquirer")
}

func TestAcquirerValidator(t *testing.T) {
	testCases := []struct {
		TestName     string
		AcquirerName string
		Err          *errors.ValidationError
	}{
		{
			"name is emtpy",
			"",
			errors.NewValidationError("acquirer name is required"),
		},
		{
			"all fields are valid",
			"Acquirer",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.TestName, func(t *testing.T) {
			err := NewAcquirer(tc.AcquirerName).Validate()
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
