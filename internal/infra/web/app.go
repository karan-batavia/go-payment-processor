package web

import (
	"github.com/sesaquecruz/go-payment-processor/internal/infra/web/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "github.com/sesaquecruz/go-payment-processor/docs"
)

func InitApp(
	paymentHandler handler.IPaymentHandler,
) *fiber.App {
	app := fiber.New()

	v1 := app.Group("/api/v1")
	{
		v1.Get("/swagger/*", swagger.HandlerDefault)

		payment := v1.Group("/payments")
		{
			payment.Post("/process", paymentHandler.ProcessPayment)
		}
	}

	return app
}
