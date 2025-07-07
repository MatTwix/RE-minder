package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		claims := token.Claims.(jwt.MapClaims)

		userID := 0
		if v, ok := claims["user_id"]; ok && v != nil {
			switch val := v.(type) {
			case float64:
				userID = int(val)
			case int:
				userID = val
			}
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token: user_id missing"})
		}

		isAdmin := false
		if v, ok := claims["is_admin"].(bool); ok {
			isAdmin = v
		} else if v, ok := claims["is_admin"].(float64); ok {
			isAdmin = v == 1
		}

		c.Locals("user_id", userID)
		c.Locals("is_admin", isAdmin)

		return c.Next()
	}
}
