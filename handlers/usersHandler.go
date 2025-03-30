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
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.TelegramId, &user.CreatedAt, &user.UpdatedAt); err != nil {
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
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.TelegramId, &user.CreatedAt, &user.UpdatedAt)

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
		"INSERT INTO users (username, email, password_hash, telegram_id) VALUES ($1, $2, $3, $4)",
		user.Username, user.Email, user.PasswordHash, user.TelegramId)
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
		SET username = $1, email = $2, password_hash = $3, telegram_id = $4, updated_at = NOW()
		WHERE id = $5`,
		user.Username, user.Email, user.PasswordHash, user.TelegramId, id)

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
