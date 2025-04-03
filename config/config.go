package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DbUser     string
	DbPort     string
	DbPassword string
	DbName     string
	ENV        string
	AppUrl     string
	ReactPort  string
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
		Port:       port,
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbPort:     os.Getenv("DB_PORT"),
		DbName:     os.Getenv("DB_NAME"),
		ENV:        os.Getenv("ENV"),
		AppUrl:     os.Getenv("APP_URL"),
		ReactPort:  os.Getenv("REACT_PORT"),
	}
}
