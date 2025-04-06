package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ENV    string
	AppUrl string
	Port   string

	DbUser     string
	DbPort     string
	DbPassword string
	DbName     string

	ReactPort string

	GithubClient       string
	GithubClientSecret string
	GithubRedirectUrl  string
	JWTSecret          string
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
		Port:   port,
		ENV:    os.Getenv("ENV"),
		AppUrl: os.Getenv("APP_URL"),

		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbPort:     os.Getenv("DB_PORT"),
		DbName:     os.Getenv("DB_NAME"),

		ReactPort: os.Getenv("REACT_PORT"),

		GithubClient:       os.Getenv("GITHUB_CLIENT"),
		GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		GithubRedirectUrl:  os.Getenv("GITHUB_REDIRECT_URL"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
	}
}
