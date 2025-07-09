package handlers

import (
	"strconv"

	"github.com/MatTwix/RE-minder/config"
	"github.com/MatTwix/RE-minder/services"
	"github.com/gofiber/fiber/v3"
)

type usersInput struct {
	Username   string `json:"username" validate:"required,min=1,max=39"`
	TelegramId *int   `json:"telegram_id,omitempty" validate:"omitempty,min=1"`
	GithubId   *int   `json:"github_id,omitempty" validate:"omitempty,min=1"`
}

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
	var input usersInput

	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Incorrect data format: " + err.Error()})
	}

	if err := config.Validator.Struct(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Validation error: " + err.Error()})
	}

	createdUser, err := services.CreateUser(c.Context(), input.Username, input.TelegramId, input.GithubId)
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

	var input usersInput

	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Incorrect data format: " + err.Error()})
	}

	if err := config.Validator.Struct(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Validation error: " + err.Error()})
	}

	updatedUser, err := services.UpdateUser(c.Context(), id, input.Username, input.TelegramId, input.GithubId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error updating user: " + err.Error()})
	}
	return c.JSON(updatedUser)
}

func SetTelegramID(c fiber.Ctx) error {
	idRaw := c.Params("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID: " + err.Error()})
	}

	var requestBody struct {
		TelegramId int `json:"telegram_id"`
	}

	if err := c.Bind().Body(&requestBody); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Incorrect data format: " + err.Error()})
	}

	if requestBody.TelegramId <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Telegram ID"})
	}

	existingUser, err := services.GetUsers(c.Context(), services.Condition{
		Field:    "id",
		Operator: services.Equal,
		Value:    id,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error while getting user: " + err.Error()})
	}

	if len(existingUser) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	singleUser := existingUser[0]

	updatedUser, err := services.UpdateUser(c.Context(), id, singleUser.Username, &requestBody.TelegramId, singleUser.GithubId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error updating user: " + err.Error()})
	}
	return c.JSON(updatedUser)
}

func SwapUserStatus(c fiber.Ctx) error {
	idRaw := c.Params("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID: " + err.Error()})
	}

	user, err := services.GetUsers(c.Context(), services.Condition{
		Field:    "id",
		Operator: services.Equal,
		Value:    id,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error while getting user: " + err.Error()})
	}

	if len(user) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	singleUser := user[0]

	isAdmin := !singleUser.IsAdmin

	updatedUser, err := services.SetUserStatus(c.Context(), id, isAdmin)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error updating user: " + err.Error()})
	}

	updatedUser.IsAdmin = isAdmin

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
