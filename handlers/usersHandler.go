package handlers

import (
	"context"

	"github.com/MatTwix/RE-minder/database"
	"github.com/MatTwix/RE-minder/models"
	"github.com/gofiber/fiber/v3"
)

func GetUsers(c fiber.Ctx) error {
	rows, err := database.DB.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error while getting users: " + err.Error()})
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.TelegramId, &user.GithubId, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error parsing data: " + err.Error()})
		}
		users = append(users, user)
	}

	return c.JSON(users)
}

func GetUser(c fiber.Ctx) error {
	id := c.Params("id")

	var user models.User
	err := database.DB.QueryRow(context.Background(),
		"SELECT * FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.GithubId, &user.TelegramId, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found: " + err.Error()})
	}

	return c.JSON(user)
}

func CreateUser(c fiber.Ctx) error {
	user := new(models.User)
	if err := c.Bind().Body(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Incorrect data format: " + err.Error()})
	}

	_, err := database.DB.Exec(context.Background(),
		"INSERT INTO users (username, telegram_id, github_id) VALUES ($1, $2, $3)",
		user.Username, user.TelegramId, user.GithubId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error creating user: " + err.Error()})
	}

	return c.JSON(user)
}

func UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")
	user := new(models.User)

	if err := c.Bind().Body(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Incorrect data format: " + err.Error()})
	}

	_, err := database.DB.Exec(context.Background(), `
		UPDATE users 
		SET username = $1, telegram_id = $2, github_id = $3, updated_at = NOW()
		WHERE id = $4`,
		user.Username, user.TelegramId, user.GithubId, id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error updating user: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User successfully updated"})
}

func DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")

	_, err := database.DB.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error deleting user: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User successfully deleted"})
}
