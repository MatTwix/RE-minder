package handlers

import (
	"github.com/MatTwix/RE-minder/services"
	"github.com/gofiber/fiber/v3"
)

func GetMe(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	user, err := services.GetUsers(c.Context(), services.Condition{
		Field:    "id",
		Operator: services.Equal,
		Value:    userID,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error while getting user: " + err.Error()})
	}

	if len(user) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	singleUser := user[0]

	return c.JSON(singleUser)
}
