package routes

import (
	"github.com/MatTwix/RE-minder/handlers"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {
	users := app.Group("/users")
	users.Get("/", handlers.GetUsers)
}
