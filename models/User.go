package models

import "time"

type User struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	TelegramId *int      `json:"telegram_id,omitempty"`
	GithubId   *int      `json:"github_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
