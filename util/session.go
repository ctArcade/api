package util

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func GetSession(c *fiber.Ctx) jwt.MapClaims {
	userLocal := c.Locals("user")

	if userLocal == nil {
		return nil
	}

	loggedU := userLocal.(*jwt.Token)
	claims := loggedU.Claims.(jwt.MapClaims)

	return claims
}
