package acquirer

import (
	"context"
	"net/http"

	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
)

type IAcquirer interface {
	Name() string
	RequestBuilder(context.Context, *entity.Transaction) (*http.Request, error)
	ResponseExtractor(*http.Response) (*entity.Payment, error)
}
