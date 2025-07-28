package services

import (
	"context"
	"errors"

	"github.com/MatTwix/RE-minder/database"
	"github.com/MatTwix/RE-minder/models"
	"github.com/jackc/pgx/v5"
)

func GetUserNotificationSettings(userID int) (models.NotificationSettings, error) {
	ctx := context.Background()
	var settings models.NotificationSettings

	err := database.DB.QueryRow(ctx, `
		SELECT id, user_id, telegram_notification, discord_notification, google_notification, created_at, updated_at
		FROM notifications_settings
		WHERE user_id = $1`, userID).Scan(
		&settings.ID,
		&settings.UserID,
		&settings.TelegramNotification,
		&settings.DiscordNotification,
		&settings.GoogleNotification,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			err := database.DB.QueryRow(context.Background(), `
				INSERT INTO notifications_settings (user_id)
				VALUES ($1)
				RETURNING id`, userID).Scan(&settings.ID)
			if err != nil {
				return settings, err
			}
			settings.UserID = userID
		} else {
			return settings, err
		}
	}

	return settings, nil
}

func UpdateUserNotificationSettings(userID int, telegram, discord, google bool) (models.NotificationSettings, error) {
	ctx := context.Background()
	settings, err := GetUserNotificationSettings(userID)
	if err != nil {
		return settings, errors.New("Error while getting user notification settings: " + err.Error())
	}

	err = database.DB.QueryRow(ctx, `
		UPDATE notifications_settings
		SET telegram_notification = $1, discord_notification = $2, google_notification = $3, updated_at = NOW()
		WHERE user_id = $4
		RETURNING telegram_notification, discord_notification, google_notification, updated_at`,
		telegram, discord, google, userID).Scan(
		&settings.TelegramNotification,
		&settings.DiscordNotification,
		&settings.GoogleNotification,
		&settings.UpdatedAt,
	)

	if err != nil {
		return settings, err
	}

	return settings, nil
}
