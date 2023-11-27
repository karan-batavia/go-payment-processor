package router

import (
	"github.com/sesaquecruz/go-payment-processor/internal/infra/web/handler"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(
	app *fiber.App,
	paymentHandler handler.IPaymentHandler,
) {
	v1 := app.Group("/api/v1")
	{
		payment := v1.Group("/payment")
		{
			payment.Post("/process", paymentHandler.ProcessPayment)
		}
	}
}
