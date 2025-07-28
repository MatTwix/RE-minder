package migrations

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func СhangeVkFieldsToGoogleNotificationsTable(DB *pgxpool.Pool) {
	ctx := context.Background()
	tx, err := DB.Begin(ctx)
	if err != nil {
		log.Fatal("Error starting transaction for renaming vk_notification: " + err.Error())
		return
	}
	defer tx.Rollback(ctx)

	var columnExists bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_name = 'notifications_settings' AND column_name = 'vk_notification'
		);
		`).Scan(&columnExists)

	if err != nil {
		log.Fatal("Error checking for vk_notification column: " + err.Error())
		return
	}

	if columnExists {
		_, err = tx.Exec(ctx, `
			ALTER TABLE notifications_settings
			RENAME COLUMN vk_notification TO google_notification;
		`)
		if err != nil {
			log.Fatal("Error renaming column vk_notificaion: " + err.Error())
			return
		}

		err = tx.Commit(ctx)
		if err != nil {
			log.Fatal("Error comitting transaction for renaming column: " + err.Error())
			return
		}
		log.Println("Successfully renamed column 'vk_notification'!")
	}
}
