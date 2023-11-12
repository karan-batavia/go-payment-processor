package service

import (
	"context"

	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
)

type PaymentService interface {
	ProcessTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Payment, error)
}
