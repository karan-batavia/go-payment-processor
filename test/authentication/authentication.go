package authentication

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func App() *fiber.App {
	app := fiber.New()

	app.Get("/token", func(c *fiber.Ctx) error {
		token, err := GetAuthToken()
		if err != nil {
			slog.Error(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": token})
	})

	return app
}

func GetAuthToken() (string, error) {
	claims := jwt.MapClaims{
		"service-id": uuid.NewString(),
		"exp":        time.Now().Add(5 * time.Minute).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token, err := jwtToken.SignedString(privateKey)

	return token, err
}
