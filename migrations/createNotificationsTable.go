package migrations

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateNotificationsTable(DB *pgxpool.Pool) {
	ctx := context.Background()
	tx, err := DB.Begin(ctx)
	if err != nil {
		log.Fatal("Error creating notifications table: " + err.Error())
	}
	defer tx.Rollback(ctx)

	var tableExists bool
	err = tx.QueryRow(ctx,
		"SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'notifications');").
		Scan(&tableExists)
	if err != nil {
		log.Fatal("Error while checking notifications table: " + err.Error())
	}

	if !tableExists {
		_, err = tx.Exec(ctx, `
			CREATE TABLE notifications (
				id SERIAL PRIMARY KEY,
				user_id INTEGER NOT NULL,
				telegram_notification BOOLEAN NOT NULL DEFAULT FALSE,
				discord_notification BOOLEAN NOT NULL DEFAULT FALSE,
				vk_notification BOOLEAN NOT NULL DEFAULT FALSE,
				created_at TIMESTAMP NOT NULL DEFAULT NOW(),
				updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
				CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
			);
		`)
		if err != nil {
			log.Fatal("Error creating notifications table: ", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			log.Fatal("Error committing transaction: ", err)
		}

		log.Println("Notifications table successfully created!")
	} else {
		tx.Rollback(ctx)
	}
}
