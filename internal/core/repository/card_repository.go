package repository

import (
	"context"

	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
)

type ICardRepository interface {
	FindCard(ctx context.Context, cardToken string) (*entity.Card, error)
}
