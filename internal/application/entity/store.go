package entity

import "github.com/sesaquecruz/go-payment-processor/internal/application/errors"

const (
	StoreIdentificationIsRequiredErr = errors.Validation("store identification is requred")
	StoreAddressIsRequiredErr        = errors.Validation("store address is required")
	StoreCepIsRequiredErr            = errors.Validation("store cep is required")
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
	if s.Identification == "" {
		return StoreIdentificationIsRequiredErr
	}

	if s.Address == "" {
		return StoreAddressIsRequiredErr
	}

	if s.Cep == "" {
		return StoreCepIsRequiredErr
	}

	return nil
}
