package repository

import (
	"context"

	app_errors "github.com/sesaquecruz/go-payment-processor/internal/application/errors"

	"github.com/sesaquecruz/go-payment-processor/internal/application/entity"
)

const (
	ErrorCardTokenIsInvalid = app_errors.Error("card token is invalid")
)

type Card interface {
	Find(ctx context.Context, cardToken string) (*entity.Card, error)
}
