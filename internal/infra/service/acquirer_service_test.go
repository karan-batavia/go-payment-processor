package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	app_errors "github.com/sesaquecruz/go-payment-processor/internal/application/errors"

	"github.com/sesaquecruz/go-payment-processor/acquirer"
	"github.com/sesaquecruz/go-payment-processor/internal/application/entity"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AcquirerServiceTestSuite struct {
	suite.Suite
	ctx             context.Context
	acquirerApp     *fiber.App
	acquirerService *AcquirerService
}

func (s *AcquirerServiceTestSuite) SetupSuite() {
	acquirerUrl := "http://127.0.0.1:6062"
	acquirerApp := acquirer.App()

	go func() {
		acquirerApp.Listen(":6062")
	}()

	requestBuilder := func(url string) AcquirerRequestBuilder {
		return func(ctx context.Context, transaction *entity.Transaction) (*http.Request, error) {
			body, err := json.Marshal(transaction)
			if err != nil {
				return nil, app_errors.NewInternal(err)
			}

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			if err != nil {
				return nil, app_errors.NewInternal(err)
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
			return nil, app_errors.NewInternal(err)
		}

		var responseData Response
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			return nil, app_errors.NewInternal(err)
		}

		if response.StatusCode != http.StatusOK {
			return nil, app_errors.NewAcquirer(responseData.Code, app_errors.Error(responseData.Message))
		}

		payment := entity.NewPayment(responseData.Message, entity.PaymentStatusPaid)
		return payment, nil
	}

	httpClient := &http.Client{}

	stone := NewAcquirer("stone", requestBuilder(acquirerUrl+"/stone"), responseExtractor)
	cielo := NewAcquirer("cielo", requestBuilder(acquirerUrl+"/cielo"), responseExtractor)
	rede := NewAcquirer("rede", requestBuilder(acquirerUrl+"/rede"), responseExtractor)

	acquirerService := NewAcquirerService(
		WithHttpClient(httpClient),
		WithAcquirer(stone),
		WithAcquirer(cielo),
		WithAcquirer(rede),
	)

	s.ctx = context.Background()
	s.acquirerApp = acquirerApp
	s.acquirerService = acquirerService
}

func (s *AcquirerServiceTestSuite) TearDownSuite() {
	if err := s.acquirerApp.Shutdown(); err != nil {
		s.FailNow(err.Error())
	}
}

func (s *AcquirerServiceTestSuite) TestAcquirerStone() {
	acquirer := "stone"

	s.T().Run("process the transaction successfully", func(t *testing.T) {
		transaction := createTransaction(acquirer, 100)
		payment, err := s.acquirerService.Process(s.ctx, transaction)
		require.Nil(t, err)
		assert.NotEmpty(t, payment.Id)
		assert.Equal(t, entity.PaymentStatusPaid, payment.Status)
	})

	s.T().Run("fails to process the transaction", func(t *testing.T) {
		transaction := createTransaction(acquirer, 101)
		_, err := s.acquirerService.Process(s.ctx, transaction)
		require.NotNil(t, err)

		var e *app_errors.Acquirer
		require.ErrorAs(t, err, &e)
		assert.Equal(t, http.StatusUnprocessableEntity, e.Code)
		assert.Equal(t, "the maximum purchase value should not exceed 100", e.Err.Error())
	})
}

func (s *AcquirerServiceTestSuite) TestAcquirerCielo() {
	acquirer := "cielo"

	s.T().Run("process the transaction successfully", func(t *testing.T) {
		transaction := createTransaction(acquirer, 500)
		payment, err := s.acquirerService.Process(s.ctx, transaction)
		require.Nil(t, err)
		assert.NotEmpty(t, payment.Id)
		assert.Equal(t, entity.PaymentStatusPaid, payment.Status)
	})

	s.T().Run("fails to process the transaction", func(t *testing.T) {
		transaction := createTransaction(acquirer, 501)
		_, err := s.acquirerService.Process(s.ctx, transaction)
		require.NotNil(t, err)

		var e *app_errors.Acquirer
		require.ErrorAs(t, err, &e)
		assert.Equal(t, http.StatusUnprocessableEntity, e.Code)
		assert.Equal(t, "the maximum purchase value should not exceed 500", e.Err.Error())
	})
}

func (s *AcquirerServiceTestSuite) TestAcquirerRede() {
	acquirer := "rede"

	s.T().Run("process the transaction successfully", func(t *testing.T) {
		transaction := createTransaction(acquirer, 1000)
		payment, err := s.acquirerService.Process(s.ctx, transaction)
		require.Nil(t, err)
		assert.NotEmpty(t, payment.Id)
		assert.Equal(t, entity.PaymentStatusPaid, payment.Status)
	})

	s.T().Run("fails to process the transaction", func(t *testing.T) {
		transaction := createTransaction(acquirer, 1001)
		_, err := s.acquirerService.Process(s.ctx, transaction)
		require.NotNil(t, err)

		var e *app_errors.Acquirer
		require.ErrorAs(t, err, &e)
		assert.Equal(t, http.StatusUnprocessableEntity, e.Code)
		assert.Equal(t, "the maximum purchase value should not exceed 1000", e.Err.Error())
	})
}

func TestAcquirerServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AcquirerServiceTestSuite))
}

func createTransaction(acquirerName string, value float64) *entity.Transaction {
	card := entity.NewCard("Token", "Holder", "01/2030", "Brand")
	purchase := entity.NewPurchase(value, []string{"Item 1", "Item 2"}, 2)
	store := entity.NewStore("Identification", "Address", "Cep")
	acquirer := entity.NewAcquirer(acquirerName)

	return entity.NewTransaction(*card, *purchase, *store, *acquirer)
}
