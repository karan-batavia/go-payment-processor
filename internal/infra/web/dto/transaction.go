package dto

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Card struct {
	Token string `json:"token" validate:"required"`
}

type Purchase struct {
	Value        float64  `json:"value"        validate:"required"`
	Itens        []string `json:"items"        validate:"required"`
	Installments int      `json:"installments" validate:"required"`
}

type Store struct {
	Identification string `json:"identification" validate:"required"`
	Address        string `json:"address"        validate:"required"`
	Cep            string `json:"cep"            validate:"required"`
}

type Acquirer struct {
	Name string `json:"name" validate:"required"`
}

type Transaction struct {
	Card     Card     `json:"card"`
	Purchase Purchase `json:"purchase"`
	Store    Store    `json:"store"`
	Acquirer Acquirer `json:"acquirer"`
}

func (t *Transaction) Validate() error {
	err := validate.Struct(t)
	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)
	msgs := make([]string, 0, len(errs))

	for _, oe := range errs {
		msg := fmt.Sprintf("%s is required", strings.ToLower(strings.ReplaceAll(oe.Namespace(), ".", " ")))
		msgs = append(msgs, msg)
	}

	return NewError(msgs...)
}
