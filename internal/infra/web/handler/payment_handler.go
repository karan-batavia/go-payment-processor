package handler

import (
	"github.com/sesaquecruz/go-payment-processor/internal/core/usecase"
	"github.com/sesaquecruz/go-payment-processor/internal/infra/web/dto"
	"github.com/sesaquecruz/go-payment-processor/internal/infra/web/errors"

	"github.com/gofiber/fiber/v2"
)

type IPaymentHandler interface {
	ProcessPayment(c *fiber.Ctx) error
}

type PaymentHandler struct {
	processPayment usecase.IProcessPayment
}

func NewPaymentHandler(processPayment usecase.IProcessPayment) *PaymentHandler {
	return &PaymentHandler{
		processPayment: processPayment,
	}
}

func (h *PaymentHandler) ProcessPayment(c *fiber.Ctx) error {
	transaction := dto.Transaction{}
	err := c.BodyParser(&transaction)
	if err != nil {
		return errors.NewHttpError(c, err)
	}

	err = transaction.Validate()
	if err != nil {
		return errors.NewHttpError(c, err)
	}

	input := usecase.ProcessPaymentInput{
		CardToken:            transaction.Card.Token,
		PurchaseValue:        transaction.Purchase.Value,
		PurchaseItems:        transaction.Purchase.Itens,
		PurchaseInstallments: transaction.Purchase.Installments,
		StoreIdentification:  transaction.Store.Identification,
		StoreAddress:         transaction.Store.Address,
		StoreCep:             transaction.Store.Cep,
		AcquirerName:         transaction.Acquirer.Name,
	}

	output, err := h.processPayment.Execute(c.Context(), &input)
	if err != nil {
		return errors.NewHttpError(c, err)
	}

	payment := dto.NewPayment(output.PaymentId)
	return c.JSON(payment)
}
