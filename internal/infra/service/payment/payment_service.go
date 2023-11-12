package payment

import (
	"context"
	"net/http"

	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
)

const (
	errorAcquirerIsInvalid = errors.Error("acquirer is invalid")
)

type PaymentService struct {
	httpClient *http.Client
	acquirers  map[string]*Acquirer
}

func NewPaymentService(options ...option) *PaymentService {
	service := &PaymentService{
		httpClient: &http.Client{},
		acquirers:  make(map[string]*Acquirer),
	}

	for _, option := range options {
		option(service)
	}

	return service
}

func (s *PaymentService) ProcessTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Payment, error) {
	acquirer, ok := s.acquirers[transaction.Acquirer.Name]
	if !ok {
		return nil, errors.NewNotFoundError(errorAcquirerIsInvalid)
	}

	request, err := acquirer.requestBuilder(ctx, transaction)
	if err != nil {
		return nil, err
	}

	response, err := s.httpClient.Do(request)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	defer response.Body.Close()
	payment, err := acquirer.responseExtractor(response)

	return payment, err
}
