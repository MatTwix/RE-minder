package handlers

import (
	"strconv"

	"github.com/MatTwix/RE-minder/services"
	"github.com/gofiber/fiber/v3"
)

type notificationSettingsInput struct {
	TelegramNotification bool `json:"telegram_notification"`
	DiscordNotification  bool `json:"discord_notification"`
	GoogleNotification   bool `json:"google_notification"`
}

func GetUserNotificationSettings(c fiber.Ctx) error {
	userIDRaw := c.Params("id")
	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID: " + err.Error()})
	}

	settings, err := services.GetUserNotificationSettings(userID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error while getting notification settings: " + err.Error()})
	}
	return c.JSON(settings)
}

func UpdateUserNotificationSettings(c fiber.Ctx) error {
	userIDRaw := c.Params("id")
	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID: " + err.Error()})
	}

	var input notificationSettingsInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input: " + err.Error()})
	}

	settings, err := services.UpdateUserNotificationSettings(userID, input.TelegramNotification, input.DiscordNotification, input.GoogleNotification)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error while updating notification settings: " + err.Error()})
	}

	return c.JSON(settings)
}
