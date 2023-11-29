package service

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/sesaquecruz/go-payment-processor/internal/acquirer"
	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
	core_errors "github.com/sesaquecruz/go-payment-processor/internal/core/errors"
)

type PaymentOption func(*PaymentService)

func PaymentWithHttpClient(httpClient *http.Client) PaymentOption {
	return func(s *PaymentService) {
		s.httpClient = httpClient
	}
}

func PaymentWithAcquirer(acquirer acquirer.IAcquirer) PaymentOption {
	return func(s *PaymentService) {
		s.acquirers[acquirer.Name()] = acquirer
	}
}

type PaymentService struct {
	httpClient *http.Client
	acquirers  map[string]acquirer.IAcquirer
}

func NewPaymentService(options ...PaymentOption) *PaymentService {
	service := &PaymentService{
		httpClient: &http.Client{},
		acquirers:  make(map[string]acquirer.IAcquirer),
	}

	for _, option := range options {
		option(service)
	}

	return service
}

func (s *PaymentService) ProcessTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Payment, error) {
	acquirer, ok := s.acquirers[transaction.Acquirer.Name]
	if !ok {
		return nil, core_errors.NewNotFoundError("acquirer is invalid")
	}

	request, err := acquirer.RequestBuilder(ctx, transaction)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	response, err := s.httpClient.Do(request)
	if err != nil {
		slog.Error(err.Error())
		return nil, core_errors.NewInternalError(err)
	}

	defer response.Body.Close()
	payment, err := acquirer.ResponseExtractor(response)
	if err != nil {
		slog.Error(err.Error())
	}

	return payment, err
}
