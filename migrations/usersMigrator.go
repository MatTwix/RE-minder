package migrations

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUsersTable(DB *pgxpool.Pool) {
	ctx := context.Background()
	tx, err := DB.Begin(ctx)
	if err != nil {
		log.Fatal("Error creating users table: ", err)
	}
	defer tx.Rollback(ctx)

	var tableExists bool
	err = tx.QueryRow(ctx,
		"SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users');").
		Scan(&tableExists)
	if err != nil {
		log.Fatal("Error while checking users table: ", err)
	}

	if !tableExists {
		_, err = tx.Exec(ctx, `
			CREATE TABLE users (
				id SERIAL PRIMARY KEY,
				username TEXT UNIQUE NOT NULL,
				telegram_id BIGINT UNIQUE,
				github_id BIGINT UNIQUE NOT NULL,
				created_at TIMESTAMP DEFAULT NOW(),
				updated_at TIMESTAMP DEFAULT NOW()
			);
		`)
		if err != nil {
			log.Fatal("Error creating users table: ", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			log.Fatal("Error commiting transaction: ", err)
		}

		log.Println("Users table successfulsly created!")
	} else {
		tx.Rollback(ctx)
	}
}
