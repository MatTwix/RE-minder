package handlers

import (
	"strconv"

	"github.com/MatTwix/RE-minder/models"
	"github.com/MatTwix/RE-minder/services"
	"github.com/gofiber/fiber/v3"
)

func GetUsers(c fiber.Ctx) error {
	users, err := services.GetUsers(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error while getting users: " + err.Error()})
	}
	return c.JSON(users)
}

func GetUser(c fiber.Ctx) error {
	user, err := services.GetUsers(c.Context(), services.Condition{
		Field:    "id",
		Operator: services.Equal,
		Value:    c.Params("id"),
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

func CreateUser(c fiber.Ctx) error {
	user := new(models.User)
	if err := c.Bind().Body(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Incorrect data format: " + err.Error()})
	}

	createdUser, err := services.CreateUser(c.Context(), user.Username, user.TelegramId, user.GithubId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error creating user: " + err.Error()})
	}

	return c.JSON(createdUser)
}

func UpdateUser(c fiber.Ctx) error {
	idRaw := c.Params("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID: " + err.Error()})
	}

	user := new(models.User)

	if err := c.Bind().Body(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Incorrect data format: " + err.Error()})
	}

	updatedUser, err := services.UpdateUser(c.Context(), id, user.Username, user.TelegramId, user.GithubId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error updating user: " + err.Error()})
	}
	return c.JSON(updatedUser)
}

func DeleteUser(c fiber.Ctx) error {
	idRaw := c.Params("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID: " + err.Error()})
	}
	if err := services.DeleteUser(c.Context(), id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error deleting user: " + err.Error()})
	}
	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}
