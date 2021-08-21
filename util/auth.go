package util

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func Auth(jwtSecret string) func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtSecret),
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			return Error(c, fiber.StatusUnauthorized, e)
		},
	})
}
