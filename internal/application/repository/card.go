package repository

import (
	"context"

	"github.com/sesaquecruz/go-payment-processor/internal/application/entity"
)

type Card interface {
	Find(ctx context.Context, cardToken string) (*entity.Card, error)
}
