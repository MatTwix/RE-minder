package routes

import (
	"github.com/MatTwix/RE-minder/handlers"
	"github.com/MatTwix/RE-minder/middleware"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {
	auth := app.Group("/auth")

	github := auth.Group("/github")
	github.Get("/", middleware.RedirectToGithub)

	app.Get("/auth/github/callback", middleware.GithubCallback)

	github.Get("/callback", middleware.GithubCallback)

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
