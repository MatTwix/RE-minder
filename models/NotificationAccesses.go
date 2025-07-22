package models

import "time"

type NotificationSettings struct {
	ID                   int       `json:"id"`
	UserID               int       `json:"user_id"`
	TelegramNotification bool      `json:"telegram_notification"`
	DiscordNotification  bool      `json:"discord_notification"`
	VKNotification       bool      `json:"vk_notification"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
