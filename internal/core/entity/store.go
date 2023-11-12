package entity

import (
	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
)

const (
	ErrorStoreIdentificationIsRequired = errors.Error("store identification is requred")
	ErrorStoreAddressIsRequired        = errors.Error("store address is required")
	ErrorStoreCepIsRequired            = errors.Error("store cep is required")
)

type Store struct {
	Identification string `json:"identification"`
	Address        string `json:"address"`
	Cep            string `json:"cep"`
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
		return errors.NewValidationError(errs...)
	}

	return nil
}
