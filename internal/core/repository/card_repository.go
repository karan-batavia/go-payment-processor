package repository

import (
	"context"

	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
)

type CardRepository interface {
	FindCard(ctx context.Context, cardToken string) (*entity.Card, error)
}
