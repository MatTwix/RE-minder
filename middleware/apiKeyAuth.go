package middleware

import "github.com/gofiber/fiber/v3"

func APIKeyMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" || apiKey != cfg.InternalApiKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or missing API key"})
		}

		return c.Next()
	}
}
