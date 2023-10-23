package service

import (
	"context"
	"net/http"

	app_errors "github.com/sesaquecruz/go-payment-processor/internal/application/errors"

	"github.com/sesaquecruz/go-payment-processor/internal/application/entity"
)

const (
	errorAcquirerIsInvalid = app_errors.Error("acquirer is invalid")
)

type AcquirerService struct {
	httpClient *http.Client
	acquirers  map[string]*Acquirer
}

func NewAcquirerService(options ...option) *AcquirerService {
	service := &AcquirerService{
		httpClient: &http.Client{},
		acquirers:  make(map[string]*Acquirer),
	}

	for _, option := range options {
		option(service)
	}

	return service
}

func (s *AcquirerService) Process(ctx context.Context, transaction *entity.Transaction) (*entity.Payment, error) {
	acquirer, ok := s.acquirers[transaction.Acquirer.Name]
	if !ok {
		return nil, app_errors.NewNotFound(errorAcquirerIsInvalid)
	}

	request, err := acquirer.requestBuilder(ctx, transaction)
	if err != nil {
		return nil, err
	}

	response, err := s.httpClient.Do(request)
	if err != nil {
		return nil, app_errors.NewInternal(err)
	}

	defer response.Body.Close()
	payment, err := acquirer.responseExtractor(response)

	return payment, err
}
