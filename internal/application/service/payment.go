package service

import (
	"context"

	"github.com/sesaquecruz/go-payment-processor/internal/application/entity"
)

type Payment interface {
	Process(ctx context.Context, transaction entity.Transaction) (*entity.Payment, error)
}
