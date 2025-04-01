package models

import "time"

type Habit struct {
	ID          int       `json:"id"`
	UserId      int       `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Frequency   string    `json:"frequency"`
	RemindTime  string    `json:"remind_time"`
	Timezone    string    `json:"timezone"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
