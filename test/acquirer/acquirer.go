package acquirer

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

type (
	transaction struct {
		CardToken            string   `json:"card_token"            validate:"required"`
		CardHolder           string   `json:"card_holder"           validate:"required"`
		CardExpiration       string   `json:"card_expiration"       validate:"required"`
		CardBrand            string   `json:"card_brand"            validate:"required"`
		PurchaseValue        float64  `json:"purchase_value"        validate:"required"`
		PurchaseItems        []string `json:"purchase_items"        validate:"required"`
		PurchaseInstallments int      `json:"purchase_installments" validate:"required"`
		StoreIdentification  string   `json:"store_identification"  validate:"required"`
		StoreAddress         string   `json:"store_address"         validate:"required"`
		StoreCep             string   `json:"store_cep"             validate:"required"`
	}

	response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
)

func App() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())

	app.Use(func(c *fiber.Ctx) error {
		if string(c.Request().Header.ContentType()) != fiber.MIMEApplicationJSON {
			return c.JSON(&response{http.StatusNotAcceptable, "the data type should be 'application/json'"})
		}
		return c.Next()
	})

	app.Post("/stone", handler(func(t *transaction) error {
		if t.PurchaseValue > 100 {
			return errors.New("the maximum purchase value should not exceed 100")
		}
		return nil
	}))

	app.Post("/cielo", handler(func(t *transaction) error {
		if t.PurchaseValue > 500 {
			return errors.New("the maximum purchase value should not exceed 500")
		}
		return nil
	}))

	app.Post("/rede", handler(func(t *transaction) error {
		if t.PurchaseValue > 1000 {
			return errors.New("the maximum purchase value should not exceed 1000")
		}
		return nil
	}))

	return app
}

func handler(process func(t *transaction) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var t transaction

		err := c.BodyParser(&t)
		if err == nil {
			err = validate.Struct(t)
		}

		if err != nil {
			slog.Error(err.Error())
			c.Status(http.StatusBadRequest)
			return c.JSON(&response{http.StatusBadRequest, "invalid request"})
		}

		err = process(&t)
		if err != nil {
			slog.Error(err.Error())
			c.Status(http.StatusUnprocessableEntity)
			return c.JSON(&response{http.StatusUnprocessableEntity, err.Error()})
		}

		return c.JSON(&response{http.StatusOK, uuid.NewString()})
	}
}
