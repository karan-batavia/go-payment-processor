package service

import (
	"context"

	"github.com/sesaquecruz/go-payment-processor/internal/application/entity"
	"github.com/sesaquecruz/go-payment-processor/internal/application/errors"
)

const (
	AcquirerIsInvalidErr          = errors.NotFound("acquirer is invalid")
	AcquirerIsUnavailableErr      = errors.Payment("acquirer is unavailable")
	TransactionWasNotProcessedErr = errors.Payment("transaction was not processed by the acquirer")
)

type Payment interface {
	Process(ctx context.Context, transaction *entity.Transaction) (*entity.Payment, error)
}
