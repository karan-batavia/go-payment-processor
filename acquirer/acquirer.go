package acquirer

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

type (
	card struct {
		Token      string `json:"token"      validate:"required"`
		Holder     string `json:"holder"     validate:"required"`
		Expiration string `json:"expiration" validate:"required"`
		Brand      string `json:"brand"      validate:"required"`
	}

	purchase struct {
		Value        float64  `json:"value"        validate:"required"`
		Items        []string `json:"items"        validate:"required"`
		Installments int      `json:"installments" validate:"required"`
	}

	store struct {
		Identification string `json:"identification" validate:"required"`
		Address        string `json:"address"        validate:"required"`
		Cep            string `json:"cep"            validate:"required"`
	}

	transaction struct {
		Card     card     `json:"card"     validate:"required"`
		Purchase purchase `json:"purchase" validate:"required"`
		Store    store    `json:"store"    validate:"required"`
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
	config := fiber.Config{
		StrictRouting: true,
	}
	app := fiber.New(config)
	app.Use(logger.New())

	app.Use(func(c *fiber.Ctx) error {
		if string(c.Request().Header.ContentType()) != fiber.MIMEApplicationJSON {
			return c.JSON(&response{http.StatusNotAcceptable, "the data type should be 'application/json'"})
		}
		return c.Next()
	})

	app.Post("/stone", handler(func(t *transaction) error {
		if t.Purchase.Value > 100 {
			return errors.New("the maximum purchase value should not exceed 100")
		}
		return nil
	}))

	app.Post("/cielo", handler(func(t *transaction) error {
		if t.Purchase.Value > 500 {
			return errors.New("the maximum purchase value should not exceed 500")
		}
		return nil
	}))

	app.Post("/rede", handler(func(t *transaction) error {
		if t.Purchase.Value > 1000 {
			return errors.New("the maximum purchase value should not exceed 1000")
		}
		return nil
	}))

	return app
}

func handler(verify func(t *transaction) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var t transaction

		err := c.BodyParser(&t)
		if err == nil {
			err = validate.Struct(t)
		}

		if err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return c.JSON(&response{http.StatusBadRequest, "invalid request"})
		}

		err = verify(&t)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusUnprocessableEntity)
			return c.JSON(&response{http.StatusUnprocessableEntity, err.Error()})
		}

		return c.JSON(&response{http.StatusOK, uuid.NewString()})
	}
}
