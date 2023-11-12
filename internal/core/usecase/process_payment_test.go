package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
	core_errors "github.com/sesaquecruz/go-payment-processor/internal/core/errors"
	"github.com/sesaquecruz/go-payment-processor/test/mocks/core/repository"
	"github.com/sesaquecruz/go-payment-processor/test/mocks/core/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessPaymentWithValidTransaction(t *testing.T) {
	ctx := context.Background()
	card := entity.NewCard("Token", "Holder", "Expiration", "Brand")

	input := ProcessPaymentInput{
		CardToken:            card.Token,
		PurchaseValue:        4.99,
		PurchaseItens:        []string{"Item 1", "Item 2"},
		PurchaseInstallments: 2,
		StoreIdentification:  "Identification",
		StoreAddress:         "Address",
		StoreCep:             "Cep",
		AcquirerName:         "Acquirer",
	}

	cardRepository := repository.NewCardRepositoryMock(t)
	cardRepository.
		EXPECT().
		FindCard(ctx, input.CardToken).
		Return(card, nil).
		Once()

	paymentService := service.NewPaymentServiceMock(t)
	paymentService.
		EXPECT().
		ProcessTransaction(ctx, mock.Anything).
		Run(func(ctx context.Context, transaction *entity.Transaction) {
			assert.Equal(t, card, &transaction.Card)
			assert.Equal(t, input.PurchaseValue, transaction.Purchase.Value)
			assert.EqualValues(t, input.PurchaseItens, transaction.Purchase.Items)
			assert.Equal(t, input.PurchaseInstallments, transaction.Purchase.Installments)
			assert.Equal(t, input.StoreIdentification, transaction.Store.Identification)
			assert.Equal(t, input.StoreAddress, transaction.Store.Address)
			assert.Equal(t, input.StoreCep, transaction.Store.Cep)
		}).
		Return(entity.NewPayment("id", entity.PaymentStatusPaid), nil).
		Once()

	processPayment := NewDefaultProcessPayment(cardRepository, paymentService)

	output, err := processPayment.Execute(ctx, input)
	assert.Nil(t, err)
	assert.Equal(t, output.PaymentStatus, entity.PaymentStatusPaid)
}

func TestProcessPaymentWithInvalidCardToken(t *testing.T) {
	ctx := context.Background()

	input := ProcessPaymentInput{
		CardToken:            "",
		PurchaseValue:        4.99,
		PurchaseItens:        []string{"Item 1", "Item 2"},
		PurchaseInstallments: 2,
		StoreIdentification:  "Identification",
		StoreAddress:         "Address",
		StoreCep:             "Cep",
		AcquirerName:         "Acquirer",
	}

	cardRepository := repository.NewCardRepositoryMock(t)
	cardRepository.
		EXPECT().
		FindCard(ctx, input.CardToken).
		Return(nil, core_errors.NewNotFoundError(errors.New("card not found"))).
		Once()

	paymentService := service.NewPaymentServiceMock(t)
	processPayment := NewDefaultProcessPayment(cardRepository, paymentService)

	output, err := processPayment.Execute(ctx, input)
	assert.Nil(t, output)

	var w *core_errors.NotFoundError
	assert.ErrorAs(t, err, &w)

	err = errors.Unwrap(w)
	assert.EqualError(t, err, "card not found")
}

func TestProcessPaymentWithInvalidPurchaseData(t *testing.T) {
	ctx := context.Background()
	card := entity.NewCard("Token", "Holder", "Expiration", "Brand")

	input := ProcessPaymentInput{
		CardToken:            card.Token,
		PurchaseValue:        0,
		PurchaseItens:        []string{""},
		PurchaseInstallments: 0,
		StoreIdentification:  "Identification",
		StoreAddress:         "Address",
		StoreCep:             "Cep",
		AcquirerName:         "Acquirer",
	}

	cardRepository := repository.NewCardRepositoryMock(t)
	cardRepository.
		EXPECT().
		FindCard(ctx, input.CardToken).
		Return(card, nil).
		Once()

	paymentService := service.NewPaymentServiceMock(t)
	processPayment := NewDefaultProcessPayment(cardRepository, paymentService)

	output, err := processPayment.Execute(ctx, input)
	assert.Nil(t, output)

	var w *core_errors.ValidationError
	assert.ErrorAs(t, err, &w)

	errs := w.Unwrap()
	assert.Equal(t, 3, len(errs))

	for i, err := range []error{
		entity.ErrorPurchaseValueIsInvalid,
		entity.ErrorPurchaseItemsIsInvalid,
		entity.ErrorPurchaseInstallmentsIsInvalid,
	} {
		assert.ErrorIs(t, err, errs[i])
	}
}

func TestProcessPaymentWithInvalidStoreData(t *testing.T) {
	ctx := context.Background()
	card := entity.NewCard("Token", "Holder", "Expiration", "Brand")

	input := ProcessPaymentInput{
		CardToken:            card.Token,
		PurchaseValue:        4.99,
		PurchaseItens:        []string{"Item 1", "Item 2"},
		PurchaseInstallments: 2,
		StoreIdentification:  "",
		StoreAddress:         "",
		StoreCep:             "",
		AcquirerName:         "Acquirer",
	}

	cardRepository := repository.NewCardRepositoryMock(t)
	cardRepository.
		EXPECT().
		FindCard(ctx, input.CardToken).
		Return(card, nil).
		Once()

	paymentService := service.NewPaymentServiceMock(t)
	processPayment := NewDefaultProcessPayment(cardRepository, paymentService)

	output, err := processPayment.Execute(ctx, input)
	assert.Nil(t, output)

	var w *core_errors.ValidationError
	assert.ErrorAs(t, err, &w)

	errs := w.Unwrap()
	assert.Equal(t, 3, len(errs))

	for i, err := range []error{
		entity.ErrorStoreIdentificationIsRequired,
		entity.ErrorStoreAddressIsRequired,
		entity.ErrorStoreCepIsRequired,
	} {
		assert.ErrorIs(t, err, errs[i])
	}
}

func TestProcessPaymentWithInvalidAcquirer(t *testing.T) {
	ctx := context.Background()
	card := entity.NewCard("Token", "Holder", "Expiration", "Brand")

	input := ProcessPaymentInput{
		CardToken:            card.Token,
		PurchaseValue:        4.99,
		PurchaseItens:        []string{"Item 1", "Item 2"},
		PurchaseInstallments: 2,
		StoreIdentification:  "Identification",
		StoreAddress:         "Address",
		StoreCep:             "Cep",
		AcquirerName:         "",
	}

	cardRepository := repository.NewCardRepositoryMock(t)
	cardRepository.
		EXPECT().
		FindCard(ctx, input.CardToken).
		Return(card, nil).
		Once()

	paymentService := service.NewPaymentServiceMock(t)
	processPayment := NewDefaultProcessPayment(cardRepository, paymentService)

	output, err := processPayment.Execute(ctx, input)
	assert.Nil(t, output)

	var w *core_errors.ValidationError
	assert.ErrorAs(t, err, &w)

	errs := w.Unwrap()
	assert.Equal(t, 1, len(errs))

	for i, err := range []error{entity.ErrorAcquirerNameIsRequired} {
		assert.ErrorIs(t, err, errs[i])
	}
}

func TestProcessPaymentWithAcquirerError(t *testing.T) {
	ctx := context.Background()
	card := entity.NewCard("Token", "Holder", "Expiration", "Brand")

	input := ProcessPaymentInput{
		CardToken:            card.Token,
		PurchaseValue:        4.99,
		PurchaseItens:        []string{"Item 1", "Item 2"},
		PurchaseInstallments: 2,
		StoreIdentification:  "Identification",
		StoreAddress:         "Address",
		StoreCep:             "Cep",
		AcquirerName:         "Acquirer",
	}

	cardRepository := repository.NewCardRepositoryMock(t)
	cardRepository.
		EXPECT().
		FindCard(ctx, input.CardToken).
		Return(card, nil).
		Once()

	paymentService := service.NewPaymentServiceMock(t)
	paymentService.
		EXPECT().
		ProcessTransaction(ctx, mock.Anything).
		Return(nil, core_errors.NewAcquirerError(503, errors.New("acquirer is unavailable"))).
		Once()

	processPayment := NewDefaultProcessPayment(cardRepository, paymentService)

	output, err := processPayment.Execute(ctx, input)
	assert.Nil(t, output)

	var w *core_errors.AcquirerError
	assert.ErrorAs(t, err, &w)

	err = errors.Unwrap(w)
	assert.Equal(t, 503, w.Code)
	assert.EqualError(t, err, "acquirer is unavailable")
}
