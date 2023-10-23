package usecase

import (
	"context"

	"github.com/sesaquecruz/go-payment-processor/internal/application/entity"
	"github.com/sesaquecruz/go-payment-processor/internal/application/repository"
	"github.com/sesaquecruz/go-payment-processor/internal/application/service"
)

type ProcessPaymentInput struct {
	CardToken            string
	PurchaseValue        float64
	PurchaseItens        []string
	PurchaseInstallments int
	StoreIdentification  string
	StoreAddress         string
	StoreCep             string
	AcquirerName         string
}

type ProcessPaymentOutput struct {
	PaymentStatus string
}

type ProcessPayment interface {
	Execute(ctx context.Context, input ProcessPaymentInput) (*ProcessPaymentOutput, error)
}

type DefaultProcessPayment struct {
	cardRepository repository.Card
	paymentService service.Payment
}

func NewDefaultProcessPayment(cardRepository repository.Card, paymentService service.Payment) *DefaultProcessPayment {
	return &DefaultProcessPayment{
		cardRepository: cardRepository,
		paymentService: paymentService,
	}
}

func (p *DefaultProcessPayment) Execute(ctx context.Context, input ProcessPaymentInput) (*ProcessPaymentOutput, error) {
	card, err := p.cardRepository.Find(ctx, input.CardToken)
	if err != nil {
		return nil, err
	}

	purchase := entity.NewPurchase(input.PurchaseValue, input.PurchaseItens, input.PurchaseInstallments)
	store := entity.NewStore(input.StoreIdentification, input.StoreAddress, input.StoreCep)
	acquirer := entity.NewAcquirer(input.AcquirerName)

	transaction := entity.NewTransaction(*card, *purchase, *store, *acquirer)
	if err = transaction.Validate(); err != nil {
		return nil, err
	}

	payment, err := p.paymentService.Process(ctx, transaction)
	if err != nil {
		return nil, err
	}

	output := &ProcessPaymentOutput{
		PaymentStatus: payment.Status,
	}

	return output, nil
}
