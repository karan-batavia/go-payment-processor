package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	core_errors "github.com/sesaquecruz/go-payment-processor/internal/core/errors"
	"github.com/sesaquecruz/go-payment-processor/internal/core/usecase"
	"github.com/sesaquecruz/go-payment-processor/internal/infra/web/dto"
	web_errors "github.com/sesaquecruz/go-payment-processor/internal/infra/web/errors"
	"github.com/sesaquecruz/go-payment-processor/internal/infra/web/handler"
	usecaseMocks "github.com/sesaquecruz/go-payment-processor/test/mocks/core/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProcessPayment(t *testing.T) {
	endpoint := "/api/v1/payment/process"

	t.Run("with valid transaction should return payment data", func(t *testing.T) {
		transaction := createTransactionDto()
		expectedPayment := &dto.Payment{Id: uuid.NewString()}

		processPaymentUsecase := usecaseMocks.NewIProcessPaymentMock(t)
		processPaymentUsecase.
			EXPECT().
			Execute(mock.Anything, mock.Anything).
			Run(func(ctx context.Context, input *usecase.ProcessPaymentInput) {
				assert.Equal(t, transaction.Card.Token, input.CardToken)
				assert.Equal(t, transaction.Purchase.Value, input.PurchaseValue)
				assert.Equal(t, transaction.Purchase.Itens, input.PurchaseItems)
				assert.Equal(t, transaction.Purchase.Installments, input.PurchaseInstallments)
				assert.Equal(t, transaction.Store.Identification, input.StoreIdentification)
				assert.Equal(t, transaction.Store.Address, input.StoreAddress)
				assert.Equal(t, transaction.Store.Cep, input.StoreCep)
				assert.Equal(t, transaction.Acquirer.Name, input.AcquirerName)
			}).
			Return(&usecase.ProcessPaymentOutput{
				PaymentId: expectedPayment.Id,
			}, nil).
			Once()

		paymentHandler := handler.NewPaymentHandler(processPaymentUsecase)
		app := InitApp(paymentHandler)

		reqBody, err := json.Marshal(&transaction)
		require.Nil(t, err)

		req := httptest.NewRequest("POST", endpoint, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		require.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		require.Nil(t, err)

		var payment *dto.Payment
		err = json.Unmarshal(resBody, &payment)
		require.Nil(t, err)
		assert.Equal(t, expectedPayment.Id, payment.Id)
	})

	t.Run("with invalid json should return status bad request", func(t *testing.T) {
		processPaymentUsecase := usecaseMocks.NewIProcessPaymentMock(t)
		paymentHandler := handler.NewPaymentHandler(processPaymentUsecase)
		app := InitApp(paymentHandler)

		req := httptest.NewRequest("POST", endpoint, nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		require.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		require.Nil(t, err)

		var httpErr *web_errors.HttpError
		err = json.Unmarshal(resBody, &httpErr)
		require.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, httpErr.Code)
		assert.Equal(t, []string{"unexpected end of JSON input"}, httpErr.Message)
	})

	t.Run("with empty transaction should return status bad request", func(t *testing.T) {
		processPaymentUsecase := usecaseMocks.NewIProcessPaymentMock(t)
		paymentHandler := handler.NewPaymentHandler(processPaymentUsecase)
		app := InitApp(paymentHandler)

		req := httptest.NewRequest("POST", endpoint, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		require.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		require.Nil(t, err)

		var httpErr *web_errors.HttpError
		err = json.Unmarshal(resBody, &httpErr)
		require.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, httpErr.Code)
		assert.Equal(t, []string{
			"transaction card token is required",
			"transaction purchase value is required",
			"transaction purchase itens is required",
			"transaction purchase installments is required",
			"transaction store identification is required",
			"transaction store address is required",
			"transaction store cep is required",
			"transaction acquirer name is required",
		}, httpErr.Message)
	})

	t.Run("with invalid transaction should return status UnprocessableEntity", func(t *testing.T) {
		transaction := createTransactionDto()

		processPaymentUsecase := usecaseMocks.NewIProcessPaymentMock(t)
		processPaymentUsecase.
			EXPECT().
			Execute(mock.Anything, mock.Anything).
			Return(nil, core_errors.NewValidationError("A validation error message")).
			Once()

		paymentHandler := handler.NewPaymentHandler(processPaymentUsecase)
		app := InitApp(paymentHandler)

		reqBody, err := json.Marshal(&transaction)
		require.Nil(t, err)

		req := httptest.NewRequest("POST", endpoint, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		require.Nil(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		require.Nil(t, err)

		var httpErr *web_errors.HttpError
		err = json.Unmarshal(resBody, &httpErr)
		require.Nil(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, httpErr.Code)
		assert.Equal(t, []string{"A validation error message"}, httpErr.Message)
	})

	t.Run("with unregistered card should return status not found", func(t *testing.T) {
		transaction := createTransactionDto()

		processPaymentUsecase := usecaseMocks.NewIProcessPaymentMock(t)
		processPaymentUsecase.
			EXPECT().
			Execute(mock.Anything, mock.Anything).
			Return(nil, core_errors.NewNotFoundError("A not found error message")).
			Once()

		paymentHandler := handler.NewPaymentHandler(processPaymentUsecase)
		app := InitApp(paymentHandler)

		reqBody, err := json.Marshal(&transaction)
		require.Nil(t, err)

		req := httptest.NewRequest("POST", endpoint, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		require.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)

		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		require.Nil(t, err)

		var httpErr *web_errors.HttpError
		err = json.Unmarshal(resBody, &httpErr)
		require.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
		assert.Equal(t, []string{"A not found error message"}, httpErr.Message)
	})

	t.Run("when occurs acquirer error should return acquirer response status", func(t *testing.T) {
		transaction := createTransactionDto()

		processPaymentUsecase := usecaseMocks.NewIProcessPaymentMock(t)
		processPaymentUsecase.
			EXPECT().
			Execute(mock.Anything, mock.Anything).
			Return(nil, core_errors.NewAcquirerError(429, "A rate limit error message")).
			Once()

		paymentHandler := handler.NewPaymentHandler(processPaymentUsecase)
		app := InitApp(paymentHandler)

		reqBody, err := json.Marshal(&transaction)
		require.Nil(t, err)

		req := httptest.NewRequest("POST", endpoint, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		require.Nil(t, err)
		assert.Equal(t, http.StatusTooManyRequests, res.StatusCode)

		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		require.Nil(t, err)

		var httpErr *web_errors.HttpError
		err = json.Unmarshal(resBody, &httpErr)
		require.Nil(t, err)
		assert.Equal(t, http.StatusTooManyRequests, httpErr.Code)
		assert.Equal(t, []string{"A rate limit error message"}, httpErr.Message)
	})

	t.Run("when occurs server error should return status internal server error", func(t *testing.T) {
		transaction := createTransactionDto()

		processPaymentUsecase := usecaseMocks.NewIProcessPaymentMock(t)
		processPaymentUsecase.
			EXPECT().
			Execute(mock.Anything, mock.Anything).
			Return(nil, core_errors.NewInternalError(errors.New("an internal error message"))).
			Once()

		paymentHandler := handler.NewPaymentHandler(processPaymentUsecase)
		app := InitApp(paymentHandler)

		reqBody, err := json.Marshal(&transaction)
		require.Nil(t, err)

		req := httptest.NewRequest("POST", endpoint, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)
		require.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		require.Nil(t, err)

		var httpErr *web_errors.HttpError
		err = json.Unmarshal(resBody, &httpErr)
		require.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
		assert.Equal(t, []string{"internal server error"}, httpErr.Message)
	})
}

func createTransactionDto() *dto.Transaction {
	return &dto.Transaction{
		Card: dto.Card{
			Token: "A card token",
		},
		Purchase: dto.Purchase{
			Value:        9.99,
			Itens:        []string{"Item 1"},
			Installments: 2,
		},
		Store: dto.Store{
			Identification: "A store identification",
			Address:        "A store address",
			Cep:            "A store cep",
		},
		Acquirer: dto.Acquirer{
			Name: "An acquirer name",
		},
	}
}
