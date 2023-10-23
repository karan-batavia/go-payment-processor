package service

import "net/http"

type option func(*AcquirerService)

func WithHttpClient(httpClient *http.Client) option {
	return func(s *AcquirerService) {
		s.httpClient = httpClient
	}
}

func WithAcquirer(acquirer *Acquirer) option {
	return func(s *AcquirerService) {
		s.acquirers[acquirer.name] = acquirer
	}
}
