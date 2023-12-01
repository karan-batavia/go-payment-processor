package handler

import (
	"github.com/sesaquecruz/go-payment-processor/internal/core/usecase"
	"github.com/sesaquecruz/go-payment-processor/internal/infra/web/dto"

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

// Process Payment godoc
//
// @Summary		Process a payment
// @Description	Process a payment transaction.
// @Tags		payments
// @Accept		json
// @Produce		json
// @Param		transaction			body			dto.Transaction		true	"Transaction"
// @Success		200	{object} 		dto.Payment
// @Failure		400	{object}		dto.HttpError
// @Failure		404	{object}		dto.HttpError
// @Failure		422	{object}		dto.HttpError
// @Router		/payments/process	[post]
func (h *PaymentHandler) ProcessPayment(c *fiber.Ctx) error {
	transaction := dto.Transaction{}
	err := c.BodyParser(&transaction)
	if err != nil {
		return dto.NewHttpError(c, err)
	}

	err = transaction.Validate()
	if err != nil {
		return dto.NewHttpError(c, err)
	}

	input := usecase.ProcessPaymentInput{
		CardToken:            transaction.CardToken,
		PurchaseValue:        transaction.PurchaseValue,
		PurchaseItems:        transaction.PurchaseItens,
		PurchaseInstallments: transaction.PurchaseInstallments,
		StoreIdentification:  transaction.StoreIdentification,
		StoreAddress:         transaction.StoreAddress,
		StoreCep:             transaction.StoreCep,
		AcquirerName:         transaction.AcquirerName,
	}

	output, err := h.processPayment.Execute(c.Context(), &input)
	if err != nil {
		return dto.NewHttpError(c, err)
	}

	payment := dto.NewPayment(output.PaymentId)
	return c.JSON(payment)
}
