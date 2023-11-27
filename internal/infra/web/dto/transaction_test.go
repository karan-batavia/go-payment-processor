package dto

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalTransaction(t *testing.T) {
	j := `
		{
			"card": {
				"token": "A Card Token"
			},
			"purchase": {
				"value": 9.99,
				"items": ["Item 1", "Item 2"],
				"installments": 3
			},
			"store": {
				"identification": "A Store Identification",
				"address": "A Store Address",
				"cep": "A Store Cep"
			},
			"acquirer": {
				"name": "An Acquirer Name"
			}
		}
	`

	transaction := Transaction{}
	err := json.Unmarshal([]byte(j), &transaction)
	require.Nil(t, err)

	assert.Equal(t, "A Card Token", transaction.Card.Token)
	assert.Equal(t, 9.99, transaction.Purchase.Value)
	assert.EqualValues(t, []string{"Item 1", "Item 2"}, transaction.Purchase.Itens)
	assert.Equal(t, 3, transaction.Purchase.Installments)
	assert.Equal(t, "A Store Identification", transaction.Store.Identification)
	assert.Equal(t, "A Store Address", transaction.Store.Address)
	assert.Equal(t, "A Store Cep", transaction.Store.Cep)
	assert.Equal(t, "An Acquirer Name", transaction.Acquirer.Name)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(transaction)
	assert.Nil(t, err)
}

func TestValidateTransaction(t *testing.T) {
	transaction := Transaction{}
	err := transaction.Validate()
	require.NotNil(t, err)

	verr := &Error{}
	require.ErrorAs(t, err, &verr)

	require.Equal(t, 8, len(verr.Messages))
	assert.Equal(t, "transaction card token is required", verr.Messages[0])
	assert.Equal(t, "transaction purchase value is required", verr.Messages[1])
	assert.Equal(t, "transaction purchase itens is required", verr.Messages[2])
	assert.Equal(t, "transaction purchase installments is required", verr.Messages[3])
	assert.Equal(t, "transaction store identification is required", verr.Messages[4])
	assert.Equal(t, "transaction store address is required", verr.Messages[5])
	assert.Equal(t, "transaction store cep is required", verr.Messages[6])
	assert.Equal(t, "transaction acquirer name is required", verr.Messages[7])
}
