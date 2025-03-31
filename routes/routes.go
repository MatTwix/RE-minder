package routes

import (
	"github.com/MatTwix/RE-minder/handlers"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	users := api.Group("/users")
	users.Get("/", handlers.GetUsers)
	users.Get("/:id", handlers.GetUser)
	users.Post("/", handlers.CreateUser)
	users.Put("/:id", handlers.UpdateUser)
	users.Delete("/:id", handlers.DeleteUser)

	habits := api.Group("/habits")
	habits.Get("/", handlers.GetHabits)
	habits.Get("/:id", handlers.GetHabit)
	habits.Post("/", handlers.CreateHabit)
	habits.Put("/:id", handlers.UpdateHabit)
	habits.Delete("/:id", handlers.DeleteHabit)
}
