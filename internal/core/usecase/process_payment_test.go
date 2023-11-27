package usecase

import (
	"context"
	"testing"

	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
	core_errors "github.com/sesaquecruz/go-payment-processor/internal/core/errors"
	"github.com/sesaquecruz/go-payment-processor/test/mocks/core/repository"
	"github.com/sesaquecruz/go-payment-processor/test/mocks/core/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProcessPaymentWithValidTransaction(t *testing.T) {
	ctx := context.Background()
	card := entity.NewCard("Token", "Holder", "Expiration", "Brand")

	input := ProcessPaymentInput{
		CardToken:            card.Token,
		PurchaseValue:        4.99,
		PurchaseItems:        []string{"Item 1", "Item 2"},
		PurchaseInstallments: 2,
		StoreIdentification:  "Identification",
		StoreAddress:         "Address",
		StoreCep:             "Cep",
		AcquirerName:         "Acquirer",
	}

	cardRepository := repository.NewICardRepositoryMock(t)
	cardRepository.
		EXPECT().
		FindCard(ctx, input.CardToken).
		Return(card, nil).
		Once()

	paymentService := service.NewIPaymentServiceMock(t)
	paymentService.
		EXPECT().
		ProcessTransaction(ctx, mock.Anything).
		Run(func(ctx context.Context, transaction *entity.Transaction) {
			assert.Equal(t, card, transaction.Card)
			assert.Equal(t, input.PurchaseValue, transaction.Purchase.Value)
			assert.EqualValues(t, input.PurchaseItems, transaction.Purchase.Items)
			assert.Equal(t, input.PurchaseInstallments, transaction.Purchase.Installments)
			assert.Equal(t, input.StoreIdentification, transaction.Store.Identification)
			assert.Equal(t, input.StoreAddress, transaction.Store.Address)
			assert.Equal(t, input.StoreCep, transaction.Store.Cep)
		}).
		Return(entity.NewPayment("id"), nil).
		Once()

	processPayment := NewProcessPayment(cardRepository, paymentService)

	output, err := processPayment.Execute(ctx, &input)
	assert.Nil(t, err)
	assert.NotEmpty(t, output.PaymentId)
}

func TestProcessPaymentWithInvalidCardToken(t *testing.T) {
	ctx := context.Background()

	input := ProcessPaymentInput{
		CardToken:            "",
		PurchaseValue:        4.99,
		PurchaseItems:        []string{"Item 1", "Item 2"},
		PurchaseInstallments: 2,
		StoreIdentification:  "Identification",
		StoreAddress:         "Address",
		StoreCep:             "Cep",
		AcquirerName:         "Acquirer",
	}

	cardRepository := repository.NewICardRepositoryMock(t)
	cardRepository.
		EXPECT().
		FindCard(ctx, input.CardToken).
		Return(nil, core_errors.NewNotFoundError("card not found")).
		Once()

	paymentService := service.NewIPaymentServiceMock(t)
	processPayment := NewProcessPayment(cardRepository, paymentService)

	output, err := processPayment.Execute(ctx, &input)
	assert.Nil(t, output)

	var w *core_errors.NotFoundError
	require.ErrorAs(t, err, &w)

	assert.Equal(t, w.Message, "card not found")
}

func TestProcessPaymentWithInvalidPurchaseData(t *testing.T) {
	ctx := context.Background()
	card := entity.NewCard("Token", "Holder", "Expiration", "Brand")

	input := ProcessPaymentInput{
		CardToken:            card.Token,
		PurchaseValue:        0,
		PurchaseItems:        []string{""},
		PurchaseInstallments: 0,
		StoreIdentification:  "Identification",
		StoreAddress:         "Address",
		StoreCep:             "Cep",
		AcquirerName:         "Acquirer",
	}

	cardRepository := repository.NewICardRepositoryMock(t)
	cardRepository.
		EXPECT().
		FindCard(ctx, input.CardToken).
		Return(card, nil).
		Once()

	paymentService := service.NewIPaymentServiceMock(t)
	processPayment := NewProcessPayment(cardRepository, paymentService)

	output, err := processPayment.Execute(ctx, &input)
	assert.Nil(t, output)

	var w *core_errors.ValidationError
	require.ErrorAs(t, err, &w)

	msgs := w.Messages
	assert.Equal(t, 3, len(msgs))

	for i, msg := range []string{
		"purchase value is invalid",
		"purchase items is invalid",
		"purchase installments is invalid",
	} {
		assert.Equal(t, msg, msgs[i])
	}
}

func TestProcessPaymentWithInvalidStoreData(t *testing.T) {
	ctx := context.Background()
	card := entity.NewCard("Token", "Holder", "Expiration", "Brand")

	input := ProcessPaymentInput{
		CardToken:            card.Token,
		PurchaseValue:        4.99,
		PurchaseItems:        []string{"Item 1", "Item 2"},
		PurchaseInstallments: 2,
		StoreIdentification:  "",
		StoreAddress:         "",
		StoreCep:             "",
		AcquirerName:         "Acquirer",
	}

	cardRepository := repository.NewICardRepositoryMock(t)
	cardRepository.
		EXPECT().
		FindCard(ctx, input.CardToken).
		Return(card, nil).
		Once()

	paymentService := service.NewIPaymentServiceMock(t)
	processPayment := NewProcessPayment(cardRepository, paymentService)

	output, err := processPayment.Execute(ctx, &input)
	assert.Nil(t, output)

	var w *core_errors.ValidationError
	require.ErrorAs(t, err, &w)

	msgs := w.Messages
	assert.Equal(t, 3, len(msgs))

	for i, msg := range []string{
		"store identification is required",
		"store address is required",
		"store cep is required",
	} {
		assert.Equal(t, msg, msgs[i])
	}
}

func TestProcessPaymentWithInvalidAcquirer(t *testing.T) {
	ctx := context.Background()
	card := entity.NewCard("Token", "Holder", "Expiration", "Brand")

	input := ProcessPaymentInput{
		CardToken:            card.Token,
		PurchaseValue:        4.99,
		PurchaseItems:        []string{"Item 1", "Item 2"},
		PurchaseInstallments: 2,
		StoreIdentification:  "Identification",
		StoreAddress:         "Address",
		StoreCep:             "Cep",
		AcquirerName:         "",
	}

	cardRepository := repository.NewICardRepositoryMock(t)
	cardRepository.
		EXPECT().
		FindCard(ctx, input.CardToken).
		Return(card, nil).
		Once()

	paymentService := service.NewIPaymentServiceMock(t)
	processPayment := NewProcessPayment(cardRepository, paymentService)

	output, err := processPayment.Execute(ctx, &input)
	assert.Nil(t, output)

	var w *core_errors.ValidationError
	require.ErrorAs(t, err, &w)

	msgs := w.Messages
	assert.Equal(t, 1, len(msgs))

	for i, msg := range []string{
		"acquirer name is required",
	} {
		assert.Equal(t, msg, msgs[i])
	}
}

func TestProcessPaymentWithAcquirerError(t *testing.T) {
	ctx := context.Background()
	card := entity.NewCard("Token", "Holder", "Expiration", "Brand")

	input := ProcessPaymentInput{
		CardToken:            card.Token,
		PurchaseValue:        4.99,
		PurchaseItems:        []string{"Item 1", "Item 2"},
		PurchaseInstallments: 2,
		StoreIdentification:  "Identification",
		StoreAddress:         "Address",
		StoreCep:             "Cep",
		AcquirerName:         "Acquirer",
	}

	cardRepository := repository.NewICardRepositoryMock(t)
	cardRepository.
		EXPECT().
		FindCard(ctx, input.CardToken).
		Return(card, nil).
		Once()

	paymentService := service.NewIPaymentServiceMock(t)
	paymentService.
		EXPECT().
		ProcessTransaction(ctx, mock.Anything).
		Return(nil, core_errors.NewAcquirerError(503, "acquirer is unavailable")).
		Once()

	processPayment := NewProcessPayment(cardRepository, paymentService)

	output, err := processPayment.Execute(ctx, &input)
	assert.Nil(t, output)

	var w *core_errors.AcquirerError
	require.ErrorAs(t, err, &w)

	assert.Equal(t, 503, w.Code)
	assert.Equal(t, "acquirer is unavailable", w.Message)
}
