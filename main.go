package main

import (
	"fmt"
	"log"

	"github.com/MatTwix/RE-minder/config"
	"github.com/MatTwix/RE-minder/database"
	"github.com/MatTwix/RE-minder/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	app := fiber.New()
	cfg := config.LoadConfig()

	database.ConnectDB()
	defer database.DB.Close()

	if cfg.ENV != "production" {
		originUrl := fmt.Sprintf("%s:%s", cfg.AppUrl, cfg.ReactPort)

		app.Use(cors.New(cors.Config{
			AllowOrigins:     []string{originUrl},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
			AllowCredentials: true,
			ExposeHeaders:    []string{"Content-Length"},
			MaxAge:           86400,
		}))
	}

	routes.SetupRoutes(app)

	if cfg.ENV == "production" {
		// TODO: react app static views from dist/ dir
	}

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
