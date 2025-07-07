package migrations

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateHabitsTable(DB *pgxpool.Pool) {
	ctx := context.Background()
	tx, err := DB.Begin(ctx)
	if err != nil {
		log.Fatal("Error creating habits table: " + err.Error())
	}
	defer tx.Rollback(ctx)

	var tableExists bool
	err = tx.QueryRow(ctx,
		"SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'habits');").
		Scan(&tableExists)
	if err != nil {
		log.Fatal("Error while checking habits table: " + err.Error())
	}

	if !tableExists {
		_, err = tx.Exec(ctx, `
			CREATE TABLE habits (
				id SERIAL PRIMARY KEY,
				user_id INTEGER NOT NULL,
				name TEXT NOT NULL,
				description TEXT DEFAULT '',
				frequency TEXT NOT NULL CHECK (frequency IN ('daily', 'weekly', 'monthly')),  
				remind_time TIME NOT NULL,
				timezone TEXT NOT NULL DEFAULT 'UTC',
				created_at TIMESTAMP NOT NULL DEFAULT NOW(),
				updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
				CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
			);
		`)
		if err != nil {
			log.Fatal("Error creating habits table: ", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			log.Fatal("Error commiting transaction: ", err)
		}

		log.Println("Habits table successfulsly created!")
	} else {
		tx.Rollback(ctx)
	}
}
