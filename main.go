package main

import (
	"log"

	"github.com/MatTwix/RE-minder/config"
	"github.com/MatTwix/RE-minder/database"
	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()
	cfg := config.LoadConfig()

	database.ConnectDB()
	defer database.DB.Close()

	if cfg.ENV != "production" {
		//TODO: CORS allows for dev mode
	}

	// TODO: make routes setup

	if cfg.ENV == "production" {
		// TODO: react app static views from dist/ dir
	}

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
