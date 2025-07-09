package migrations

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddIsAdminField(DB *pgxpool.Pool) {
	ctx := context.Background()
	tx, err := DB.Begin(ctx)
	if err != nil {
		log.Fatal("Error starting transaction to add is_admin field: ", err)
	}
	defer tx.Rollback(ctx)

	var columnExists bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'users' AND column_name = 'is_admin'
		);
	`).Scan(&columnExists)
	if err != nil {
		log.Fatal("Error checking if is_admin column exists: ", err)
	}

	if !columnExists {
		_, err = tx.Exec(ctx, `
			ALTER TABLE users ADD COLUMN is_admin BOOLEAN DEFAULT FALSE;
		`)
		if err != nil {
			log.Fatal("Error adding is_admin column: ", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			log.Fatal("Error committing transaction: ", err)
		}

		log.Println("is_admin field successfully added to users table!")
	} else {
		tx.Rollback(ctx)
	}
}
