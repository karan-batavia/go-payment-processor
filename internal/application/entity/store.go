package entity

import (
	"errors"

	app_errors "github.com/sesaquecruz/go-payment-processor/internal/application/errors"
)

var (
	ErrorStoreIdentificationIsRequired = errors.New("store identification is requred")
	ErrorStoreAddressIsRequired        = errors.New("store address is required")
	ErrorStoreCepIsRequired            = errors.New("store cep is required")
)

type Store struct {
	Identification string
	Address        string
	Cep            string
}

func NewStore(identification string, address string, cep string) *Store {
	return &Store{
		Identification: identification,
		Address:        address,
		Cep:            cep,
	}
}

func (s *Store) Validate() error {
	errs := make([]error, 0)

	if s.Identification == "" {
		errs = append(errs, ErrorStoreIdentificationIsRequired)
	}

	if s.Address == "" {
		errs = append(errs, ErrorStoreAddressIsRequired)
	}

	if s.Cep == "" {
		errs = append(errs, ErrorStoreCepIsRequired)
	}

	if len(errs) > 0 {
		return app_errors.NewValidation(errs...)
	}

	return nil
}
