package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
	"github.com/sesaquecruz/go-payment-processor/test/acquirer"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PaymentServiceTestSuite struct {
	suite.Suite
	ctx            context.Context
	acquirerApp    *fiber.App
	PaymentService *PaymentService
}

func (s *PaymentServiceTestSuite) SetupSuite() {
	acquirerUrl := "http://127.0.0.1:6062"
	acquirerApp := acquirer.App()

	go func() {
		acquirerApp.Listen(":6062")
	}()

	time.Sleep(1 * time.Second)

	requestBuilder := func(url string) AcquirerRequestBuilder {
		return func(ctx context.Context, transaction *entity.Transaction) (*http.Request, error) {
			type Transaction struct {
				CardToken            string   `json:"card_token"`
				CardHolder           string   `json:"card_holder"`
				CardExpiration       string   `json:"card_expiration"`
				CardBrand            string   `json:"card_brand"`
				PurchaseValue        float64  `json:"purchase_value"`
				PurchaseItems        []string `json:"purchase_items"`
				PurchaseInstallments int      `json:"purchase_installments"`
				StoreIdentification  string   `json:"store_identification"`
				StoreAddress         string   `json:"store_address"`
				StoreCep             string   `json:"store_cep"`
				StoreName            string   `json:"store_name"`
			}

			acquirerTransaction := Transaction{
				CardToken:            transaction.Card.Token,
				CardHolder:           transaction.Card.Holder,
				CardExpiration:       transaction.Card.Expiration,
				CardBrand:            transaction.Card.Brand,
				PurchaseValue:        transaction.Purchase.Value,
				PurchaseItems:        transaction.Purchase.Items,
				PurchaseInstallments: transaction.Purchase.Installments,
				StoreIdentification:  transaction.Store.Identification,
				StoreAddress:         transaction.Store.Address,
				StoreCep:             transaction.Store.Cep,
				StoreName:            transaction.Acquirer.Name,
			}

			body, err := json.Marshal(acquirerTransaction)
			if err != nil {
				return nil, errors.NewInternalError(err)
			}

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			if err != nil {
				return nil, errors.NewInternalError(err)
			}

			request.Header.Set("Content-Type", "application/json")
			return request, nil
		}
	}

	responseExtractor := func(response *http.Response) (*entity.Payment, error) {
		type Response struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, errors.NewInternalError(err)
		}

		var responseData Response
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			return nil, errors.NewInternalError(err)
		}

		if response.StatusCode != http.StatusOK {
			return nil, errors.NewAcquirerError(responseData.Code, responseData.Message)
		}

		payment := entity.NewPayment(responseData.Message)
		return payment, nil
	}

	httpClient := &http.Client{}

	stone := NewAcquirer("stone", requestBuilder(acquirerUrl+"/stone"), responseExtractor)
	cielo := NewAcquirer("cielo", requestBuilder(acquirerUrl+"/cielo"), responseExtractor)
	rede := NewAcquirer("rede", requestBuilder(acquirerUrl+"/rede"), responseExtractor)

	PaymentService := NewPaymentService(
		WithHttpClient(httpClient),
		WithAcquirer(stone),
		WithAcquirer(cielo),
		WithAcquirer(rede),
	)

	s.ctx = context.Background()
	s.acquirerApp = acquirerApp
	s.PaymentService = PaymentService
}

func (s *PaymentServiceTestSuite) TearDownSuite() {
	if err := s.acquirerApp.Shutdown(); err != nil {
		s.FailNow(err.Error())
	}
}

func (s *PaymentServiceTestSuite) TestAcquirerStone() {
	acquirer := "stone"

	s.T().Run("process the transaction successfully", func(t *testing.T) {
		transaction := createTransaction(acquirer, 100)
		payment, err := s.PaymentService.ProcessTransaction(s.ctx, transaction)
		require.Nil(t, err)
		assert.NotEmpty(t, payment.Id)
	})

	s.T().Run("fails to process the transaction", func(t *testing.T) {
		transaction := createTransaction(acquirer, 101)
		_, err := s.PaymentService.ProcessTransaction(s.ctx, transaction)
		require.NotNil(t, err)

		var e *errors.AcquirerError
		require.ErrorAs(t, err, &e)
		assert.Equal(t, http.StatusUnprocessableEntity, e.Code)
		assert.Equal(t, "the maximum purchase value should not exceed 100", e.Message)
	})
}

func (s *PaymentServiceTestSuite) TestAcquirerCielo() {
	acquirer := "cielo"

	s.T().Run("process the transaction successfully", func(t *testing.T) {
		transaction := createTransaction(acquirer, 500)
		payment, err := s.PaymentService.ProcessTransaction(s.ctx, transaction)
		require.Nil(t, err)
		assert.NotEmpty(t, payment.Id)
	})

	s.T().Run("fails to process the transaction", func(t *testing.T) {
		transaction := createTransaction(acquirer, 501)
		_, err := s.PaymentService.ProcessTransaction(s.ctx, transaction)
		require.NotNil(t, err)

		var e *errors.AcquirerError
		require.ErrorAs(t, err, &e)
		assert.Equal(t, http.StatusUnprocessableEntity, e.Code)
		assert.Equal(t, "the maximum purchase value should not exceed 500", e.Message)
	})
}

func (s *PaymentServiceTestSuite) TestAcquirerRede() {
	acquirer := "rede"

	s.T().Run("process the transaction successfully", func(t *testing.T) {
		transaction := createTransaction(acquirer, 1000)
		payment, err := s.PaymentService.ProcessTransaction(s.ctx, transaction)
		require.Nil(t, err)
		assert.NotEmpty(t, payment.Id)
	})

	s.T().Run("fails to process the transaction", func(t *testing.T) {
		transaction := createTransaction(acquirer, 1001)
		_, err := s.PaymentService.ProcessTransaction(s.ctx, transaction)
		require.NotNil(t, err)

		var e *errors.AcquirerError
		require.ErrorAs(t, err, &e)
		assert.Equal(t, http.StatusUnprocessableEntity, e.Code)
		assert.Equal(t, "the maximum purchase value should not exceed 1000", e.Message)
	})
}

func TestPaymentServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentServiceTestSuite))
}

func createTransaction(acquirerName string, value float64) *entity.Transaction {
	card := entity.NewCard("Token", "Holder", "01/2030", "Brand")
	purchase := entity.NewPurchase(value, []string{"Item 1", "Item 2"}, 2)
	store := entity.NewStore("Identification", "Address", "Cep")
	acquirer := entity.NewAcquirer(acquirerName)

	return entity.NewTransaction(card, purchase, store, acquirer)
}
