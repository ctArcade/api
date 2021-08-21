package util

import "github.com/gofiber/fiber/v2"

func Error(c *fiber.Ctx, status int, e error) error {
	err := fiber.Map{"error": e.Error()}

	return c.Status(status).JSON(err)
}
