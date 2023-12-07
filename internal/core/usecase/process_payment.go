package usecase

import (
	"context"

	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
	"github.com/sesaquecruz/go-payment-processor/internal/core/repository"
	"github.com/sesaquecruz/go-payment-processor/internal/core/service"
)

type ProcessPaymentInput struct {
	CardToken            string
	PurchaseValue        float64
	PurchaseItems        []string
	PurchaseInstallments int
	StoreIdentification  string
	StoreAddress         string
	StoreCep             string
	AcquirerName         string
}

type ProcessPaymentOutput struct {
	PaymentId string
}

type IProcessPayment interface {
	Execute(ctx context.Context, input *ProcessPaymentInput) (*ProcessPaymentOutput, error)
}

type ProcessPayment struct {
	cardRepository repository.ICardRepository
	paymentService service.IPaymentService
}

func NewProcessPayment(cardRepository repository.ICardRepository, paymentService service.IPaymentService) *ProcessPayment {
	return &ProcessPayment{
		cardRepository: cardRepository,
		paymentService: paymentService,
	}
}

func (p *ProcessPayment) Execute(ctx context.Context, input *ProcessPaymentInput) (*ProcessPaymentOutput, error) {
	card, err := p.cardRepository.FindCard(ctx, input.CardToken)
	if err != nil {
		return nil, err
	}

	purchase := entity.NewPurchase(input.PurchaseValue, input.PurchaseItems, input.PurchaseInstallments)
	store := entity.NewStore(input.StoreIdentification, input.StoreAddress, input.StoreCep)
	acquirer := entity.NewAcquirer(input.AcquirerName)
	transaction := entity.NewTransaction(card, purchase, store, acquirer)

	err = transaction.Validate()
	if err != nil {
		return nil, err
	}

	payment, err := p.paymentService.ProcessTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	output := &ProcessPaymentOutput{
		PaymentId: payment.Id,
	}

	return output, nil
}
