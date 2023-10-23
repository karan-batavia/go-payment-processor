package service

import (
	"context"
	"net/http"

	"github.com/sesaquecruz/go-payment-processor/internal/application/entity"
)

type AcquirerRequestBuilder func(context.Context, *entity.Transaction) (*http.Request, error)
type AcquirerResponseExtractor func(*http.Response) (*entity.Payment, error)

type Acquirer struct {
	name              string
	requestBuilder    AcquirerRequestBuilder
	responseExtractor AcquirerResponseExtractor
}

func NewAcquirer(
	name string,
	requestBuilder AcquirerRequestBuilder,
	responseExtractor AcquirerResponseExtractor,
) *Acquirer {
	return &Acquirer{
		name:              name,
		requestBuilder:    requestBuilder,
		responseExtractor: responseExtractor,
	}
}
