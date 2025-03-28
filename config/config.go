package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DbUrl string
	ENV   string
}

func LoadConfig() Config {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error while loading .env file")
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Port:  port,
		DbUrl: os.Getenv("DB_URL"),
		ENV:   os.Getenv("ENV"),
	}
}
