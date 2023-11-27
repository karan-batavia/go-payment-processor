package entity

import (
	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
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
	msgs := make([]string, 0)

	if s.Identification == "" {
		msgs = append(msgs, "store identification is required")
	}

	if s.Address == "" {
		msgs = append(msgs, "store address is required")
	}

	if s.Cep == "" {
		msgs = append(msgs, "store cep is required")
	}

	if len(msgs) > 0 {
		return errors.NewValidationError(msgs...)
	}

	return nil
}
