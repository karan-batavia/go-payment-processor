package repository

import (
	"context"

	"github.com/sesaquecruz/go-payment-processor/internal/application/entity"
	"github.com/sesaquecruz/go-payment-processor/internal/application/errors"
)

const (
	CardTokenIsInvalidErr = errors.NotFound("card token is invalid")
)

type Card interface {
	Find(ctx context.Context, cardToken string) (*entity.Card, error)
}
