package migrations

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddStartDateField(DB *pgxpool.Pool) {
	ctx := context.Background()
	tx, err := DB.Begin(ctx)
	if err != nil {
		log.Fatal("Error starting transaction to add start_date field: ", err)
	}
	defer tx.Rollback(ctx)

	var columnExists bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'habits' AND column_name = 'start_date'
		);
	`).Scan(&columnExists)
	if err != nil {
		log.Fatal("Error checking if start_date column exists: ", err)
	}

	if !columnExists {
		_, err = tx.Exec(ctx, `
			ALTER TABLE habits ADD COLUMN start_date TIMESTAMP DEFAULT NOW();
		`)
		if err != nil {
			log.Fatal("Error adding start_date column: ", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			log.Fatal("Error committing transaction: ", err)
		}

		log.Println("start_date field successfully added to habits table!")
	} else {
		tx.Rollback(ctx)
	}
}
