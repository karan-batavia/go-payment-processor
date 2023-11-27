package service

import (
	"context"

	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
)

type IPaymentService interface {
	ProcessTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Payment, error)
}
