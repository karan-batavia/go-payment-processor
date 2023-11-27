package payment

import (
	"net/http"
)

type option func(*PaymentService)

func WithHttpClient(httpClient *http.Client) option {
	return func(s *PaymentService) {
		s.httpClient = httpClient
	}
}

func WithAcquirer(acquirer *Acquirer) option {
	return func(s *PaymentService) {
		s.acquirers[acquirer.name] = acquirer
	}
}
