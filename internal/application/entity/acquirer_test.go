package entity

import (
	"testing"

	app_errors "github.com/sesaquecruz/go-payment-processor/internal/application/errors"

	"github.com/stretchr/testify/assert"
)

func TestAcquirerFactory(t *testing.T) {
	acquirer := NewAcquirer("Acquirer")
	assert.NotNil(t, acquirer)
	assert.Equal(t, acquirer.Name, "Acquirer")
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

			var v *app_errors.Validation
			assert.ErrorAs(t, err, &v)

			errs := v.Unwrap()
			assert.Equal(t, len(tc.errs), len(errs))

			for i, err := range tc.errs {
				assert.ErrorIs(t, err, errs[i])
			}
		})
	}
}
