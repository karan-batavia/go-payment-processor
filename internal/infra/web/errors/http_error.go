package errors

import (
	"encoding/json"
	"net/http"

	core_errors "github.com/sesaquecruz/go-payment-processor/internal/core/errors"
	"github.com/sesaquecruz/go-payment-processor/internal/infra/web/dto"

	"github.com/gofiber/fiber/v2"
)

type HttpError struct {
	Code    int      `json:"code"`
	Message []string `json:"message"`
}

func NewHttpError(c *fiber.Ctx, err error) error {
	httpErr := &HttpError{}

	switch t := err.(type) {
	case *json.SyntaxError:
		httpErr.Code = http.StatusBadRequest
		httpErr.Message = []string{err.Error()}
		break

	case *fiber.Error:
		httpErr.Code = http.StatusBadRequest
		httpErr.Message = []string{"Bad Request"}
		break

	case *dto.Error:
		httpErr.Code = http.StatusBadRequest
		httpErr.Message = t.Messages
		break

	case *core_errors.ValidationError:
		httpErr.Code = http.StatusUnprocessableEntity
		httpErr.Message = t.Messages
		break

	case *core_errors.NotFoundError:
		httpErr.Code = http.StatusNotFound
		httpErr.Message = []string{t.Message}
		break

	case *core_errors.AcquirerError:
		httpErr.Code = t.Code
		httpErr.Message = []string{t.Message}
		break

	default:
		httpErr.Code = http.StatusInternalServerError
		httpErr.Message = []string{"internal server error"}
	}

	return c.Status(httpErr.Code).JSON(httpErr)
}
