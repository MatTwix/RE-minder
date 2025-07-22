package migrations

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MakeUserIdFieldUniqueNotificationsSettingsTable(DB *pgxpool.Pool) {
	ctx := context.Background()
	tx, err := DB.Begin(ctx)
	if err != nil {
		log.Fatal("Error starting transaction: ", err)
	}
	defer tx.Rollback(ctx)

	var constraintExists bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage kcu
				ON tc.constraint_name = kcu.constraint_name
				AND tc.table_schema = kcu.table_schema
			WHERE tc.table_name = 'notifications_settings'
			  AND tc.constraint_type = 'UNIQUE'
			  AND kcu.column_name = 'user_id'
		);
	`).Scan(&constraintExists)

	if err != nil {
		log.Fatal("Error checking for unique constraint: ", err)
	}

	if !constraintExists {
		_, err = tx.Exec(ctx, `
			ALTER TABLE notifications_settings
			ADD CONSTRAINT notifications_settings_user_id_unique UNIQUE (user_id);
		`)
		if err != nil {
			log.Fatal("Error adding unique constraint to notification_settings table: ", err)
		}
		log.Println("Successfully added UNIQUE constraint to user_id in notifications_settings table.")
	}

	if err = tx.Commit(ctx); err != nil {
		log.Fatal("Error committing transaction: ", err)
	}
}
