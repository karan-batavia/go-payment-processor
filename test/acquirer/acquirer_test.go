package acquirer

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStoneAcquirer(t *testing.T) {
	app := App()
	url := "/stone"

	t.Run("with transaction value less than or equal to 100", func(t *testing.T) {
		reqData := createTransaction(100)
		reqBody, err := json.Marshal(reqData)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		assert.Nil(t, err)

		var resData response
		err = json.Unmarshal(resBody, &resData)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resData.Code)

		_, err = uuid.Parse(resData.Message)
		assert.Nil(t, err)
	})

	t.Run("with transaction value greater than 100", func(t *testing.T) {
		reqData := createTransaction(101)
		reqBody, err := json.Marshal(reqData)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		assert.Nil(t, err)

		var resData response
		err = json.Unmarshal(resBody, &resData)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resData.Code)
		assert.Equal(t, "the maximum purchase value should not exceed 100", resData.Message)
	})
}

func TestCieloAcquirer(t *testing.T) {
	app := App()
	url := "/cielo"

	t.Run("with transaction value less than or equal to 500", func(t *testing.T) {
		reqData := createTransaction(500)
		reqBody, err := json.Marshal(reqData)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		assert.Nil(t, err)

		var resData response
		err = json.Unmarshal(resBody, &resData)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resData.Code)

		_, err = uuid.Parse(resData.Message)
		assert.Nil(t, err)
	})

	t.Run("with transaction value greater than 500", func(t *testing.T) {
		reqData := createTransaction(501)
		reqBody, err := json.Marshal(reqData)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		assert.Nil(t, err)

		var resData response
		err = json.Unmarshal(resBody, &resData)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resData.Code)
		assert.Equal(t, "the maximum purchase value should not exceed 500", resData.Message)
	})
}

func TestRedeAcquirer(t *testing.T) {
	app := App()
	url := "/rede"

	t.Run("with transaction value less than or equal to 1000", func(t *testing.T) {
		reqData := createTransaction(1000)
		reqBody, err := json.Marshal(reqData)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		assert.Nil(t, err)

		var resData response
		err = json.Unmarshal(resBody, &resData)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, resData.Code)

		_, err = uuid.Parse(resData.Message)
		assert.Nil(t, err)
	})

	t.Run("with transaction value greater than 1000", func(t *testing.T) {
		reqData := createTransaction(1001)
		reqBody, err := json.Marshal(reqData)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		assert.Nil(t, err)

		res, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		assert.Nil(t, err)

		var resData response
		err = json.Unmarshal(resBody, &resData)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, resData.Code)
		assert.Equal(t, "the maximum purchase value should not exceed 1000", resData.Message)
	})
}

func createTransaction(value float64) *transaction {
	return &transaction{
		CardToken:            "Token",
		CardHolder:           "Holder",
		CardExpiration:       "01/2030",
		CardBrand:            "Brand",
		PurchaseValue:        value,
		PurchaseItems:        []string{"Item 1", "Item 2"},
		PurchaseInstallments: 2,
		StoreIdentification:  "Identification",
		StoreAddress:         "Address",
		StoreCep:             "Cep",
	}
}
