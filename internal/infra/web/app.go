package web

import (
	"crypto/rsa"

	"github.com/sesaquecruz/go-payment-processor/internal/infra/web/handler"

	jwtmiddleware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "github.com/sesaquecruz/go-payment-processor/docs"
)

func InitApp(
	authPublicKey *rsa.PublicKey,
	paymentHandler handler.IPaymentHandler,
) *fiber.App {
	app := fiber.New()

	v1 := app.Group("/api/v1")

	// public routes
	{
		v1.Get("/swagger/*", swagger.HandlerDefault)
	}

	v1 = v1.Use(jwtmiddleware.New(jwtmiddleware.Config{
		SigningKey: jwtmiddleware.SigningKey{
			JWTAlg: jwtmiddleware.RS256,
			Key:    authPublicKey,
		},
	}))

	// protected routes
	{
		payments := v1.Group("/payments")
		{
			payments.Post("/process", paymentHandler.ProcessPayment)
		}
	}

	return app
}
