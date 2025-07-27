package main

import (
	"fmt"
	"log"

	"github.com/MatTwix/RE-minder/config"
	"github.com/MatTwix/RE-minder/database"
	"github.com/MatTwix/RE-minder/oauth"
	"github.com/MatTwix/RE-minder/queue"
	"github.com/MatTwix/RE-minder/routes"
	"github.com/MatTwix/RE-minder/scheduler"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func main() {
	app := fiber.New()
	cfg := config.LoadConfig()

	config.InitValidator()

	database.ConnectDB()
	defer database.DB.Close()

	oauth.RegisterProvider(oauth.NewDiscordProvider())

	queue.Connect()
	scheduler.StartScheduler()

	allowedOrigins := []string{
		cfg.FrontendUrlProd,
	}

	if cfg.FrontendEnv != "production" {
		allowedOrigins = append(allowedOrigins, fmt.Sprintf("%s:%s", cfg.FrontendUrlDev, cfg.FrontendPortDev))
	}

	log.Print("Allowed origins: ", allowedOrigins)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           86400,
	}))

	if cfg.RateLimiterEnabled {
		app.Use(limiter.New(limiter.Config{
			Max:        200,
			Expiration: 60 * 1000,
			LimitReached: func(c fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Too many requests, please try again later.",
				})
			},
		}))
	}

	routes.SetupRoutes(app)

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
