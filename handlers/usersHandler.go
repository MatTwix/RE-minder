package handlers

import (
	"context"

	"github.com/MatTwix/RE-minder/database"
	"github.com/MatTwix/RE-minder/models"
	"github.com/gofiber/fiber/v3"
)

func GetUsers(c fiber.Ctx) error {
	rows, err := database.DB.Query(context.Background(), "SELECT id, name FROM users")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error while getting users: " + err.Error()})
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error parsing data: " + err.Error()})
		}
		users = append(users, user)
	}

	return c.JSON(users)
}
