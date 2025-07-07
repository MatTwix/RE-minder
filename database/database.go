package database

import (
	"context"
	"log"
	"time"

	"github.com/MatTwix/RE-minder/config"
	"github.com/MatTwix/RE-minder/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	cfg := config.LoadConfig()

	if cfg.DBUrl == "" {
		log.Fatal("There is not enough DB info in .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatal("Error trying to connect to DB: ", err)
	}

	DB = pool
	log.Print("Successful connection to DB")

	migrations.Migrate(DB)
}
