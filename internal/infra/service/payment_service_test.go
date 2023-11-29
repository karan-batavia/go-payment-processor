package service

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/sesaquecruz/go-payment-processor/internal/acquirer"
	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
	acquirer_app "github.com/sesaquecruz/go-payment-processor/test/acquirer"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PaymentServiceTestSuite struct {
	suite.Suite
	ctx            context.Context
	acquirerApp    *fiber.App
	paymentService *PaymentService
}

func (s *PaymentServiceTestSuite) SetupSuite() {
	acquirerApp := acquirer_app.App()
	acquirerUrl := "http://127.0.0.1:6062"

	go func() {
		acquirerApp.Listen(":6062")
	}()

	time.Sleep(1 * time.Second)

	paymentService := NewPaymentService(
		PaymentWithHttpClient(&http.Client{}),
		PaymentWithAcquirer(acquirer.NewCielo(acquirerUrl+"/cielo", "cielo-api-key")),
		PaymentWithAcquirer(acquirer.NewRede(acquirerUrl+"/rede", "rede-api-key")),
		PaymentWithAcquirer(acquirer.NewStone(acquirerUrl+"/stone", "stone-api-key")),
	)

	s.ctx = context.Background()
	s.acquirerApp = acquirerApp
	s.paymentService = paymentService
}

func (s *PaymentServiceTestSuite) TestAcquirerCielo() {
	acquirer := "cielo"

	s.T().Run("process the transaction successfully", func(t *testing.T) {
		transaction := createTransaction(acquirer, 100)
		payment, err := s.paymentService.ProcessTransaction(s.ctx, transaction)
		require.Nil(t, err)
		assert.NotEmpty(t, payment.Id)
	})

	s.T().Run("fails to process the transaction", func(t *testing.T) {
		transaction := createTransaction(acquirer, 101)
		_, err := s.paymentService.ProcessTransaction(s.ctx, transaction)
		require.NotNil(t, err)

		var e *errors.AcquirerError
		require.ErrorAs(t, err, &e)
		assert.Equal(t, http.StatusUnprocessableEntity, e.Code)
		assert.Equal(t, "the maximum purchase value should not exceed 100", e.Message)
	})
}

func (s *PaymentServiceTestSuite) TestAcquirerRede() {
	acquirer := "rede"

	s.T().Run("process the transaction successfully", func(t *testing.T) {
		transaction := createTransaction(acquirer, 500)
		payment, err := s.paymentService.ProcessTransaction(s.ctx, transaction)
		require.Nil(t, err)
		assert.NotEmpty(t, payment.Id)
	})

	s.T().Run("fails to process the transaction", func(t *testing.T) {
		transaction := createTransaction(acquirer, 501)
		_, err := s.paymentService.ProcessTransaction(s.ctx, transaction)
		require.NotNil(t, err)

		var e *errors.AcquirerError
		require.ErrorAs(t, err, &e)
		assert.Equal(t, http.StatusUnprocessableEntity, e.Code)
		assert.Equal(t, "the maximum purchase value should not exceed 500", e.Message)
	})
}

func (s *PaymentServiceTestSuite) TestAcquirerStone() {
	acquirer := "stone"

	s.T().Run("process the transaction successfully", func(t *testing.T) {
		transaction := createTransaction(acquirer, 1000)
		payment, err := s.paymentService.ProcessTransaction(s.ctx, transaction)
		require.Nil(t, err)
		assert.NotEmpty(t, payment.Id)
	})

	s.T().Run("fails to process the transaction", func(t *testing.T) {
		transaction := createTransaction(acquirer, 1001)
		_, err := s.paymentService.ProcessTransaction(s.ctx, transaction)
		require.NotNil(t, err)

		var e *errors.AcquirerError
		require.ErrorAs(t, err, &e)
		assert.Equal(t, http.StatusUnprocessableEntity, e.Code)
		assert.Equal(t, "the maximum purchase value should not exceed 1000", e.Message)
	})
}

func (s *PaymentServiceTestSuite) TearDownSuite() {
	if err := s.acquirerApp.Shutdown(); err != nil {
		s.FailNow(err.Error())
	}
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
