package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ENV                string
	Port               string
	RateLimiterEnabled bool

	DBUrl string

	FrontendEnv string

	FrontendUrlDev  string
	FrontendPortDev string

	FrontendUrlProd  string
	FrontendPortProd string

	GithubClient       string
	GithubClientSecret string
	JWTSecret          string
	InternalApiKey     string

	DiscordClientID     string
	DiscordClientSecret string
	VKClientID          string
	VKClientSecret      string

	DiscordRedirectUrl string
	VKRedirectUrl      string

	BotsApiUrl string

	RabbitMQUrl string
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
		Port:               port,
		ENV:                os.Getenv("ENV"),
		RateLimiterEnabled: os.Getenv("RATE_LIMITER_ENABLED") == "true",

		DBUrl: os.Getenv("DB_URL"),

		FrontendEnv: os.Getenv("FRONTEND_ENV"),

		FrontendUrlDev:  os.Getenv("FRONTEND_URL_DEV"),
		FrontendPortDev: os.Getenv("FRONTEND_PORT_DEV"),

		FrontendUrlProd:  os.Getenv("FRONTEND_URL_PROD"),
		FrontendPortProd: os.Getenv("FRONTEND_PORT_PROD"),

		GithubClient:       os.Getenv("GITHUB_CLIENT"),
		GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		InternalApiKey:     os.Getenv("INTERNAL_API_KEY"),

		DiscordClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		DiscordClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		VKClientID:          os.Getenv("VK_CLIENT_ID"),
		VKClientSecret:      os.Getenv("VK_CLIENT_SECRET"),

		DiscordRedirectUrl: os.Getenv("DISCORD_REDIRECT_URL"),
		VKRedirectUrl:      os.Getenv("VK_REDIRECT_URL"),

		BotsApiUrl: os.Getenv("BOTS_API_URL"),

		RabbitMQUrl: os.Getenv("RABBITMQ_URL"),
	}
}
