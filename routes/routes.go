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
	github.Get("/callback", middleware.GithubCallback)

	api := app.Group("/api")

	api.Get("/me", handlers.GetMe, middleware.JWTMiddleware())

	users := api.Group("/users")
	users.Get("/", handlers.GetUsers, middleware.JWTMiddleware(), middleware.IsAdmin())
	users.Get("/:id", handlers.GetUser, middleware.JWTMiddleware(), middleware.IsSelfOrAdmin())
	users.Post("/", handlers.CreateUser, middleware.JWTMiddleware(), middleware.IsAdmin())
	users.Put("/:id", handlers.UpdateUser, middleware.JWTMiddleware(), middleware.IsAdmin())
	users.Patch("/:id/telegram_id", handlers.SetTelegramID, middleware.JWTMiddleware(), middleware.IsSelfOrAdmin())
	users.Patch("/:id/is_admin", handlers.SwapUserStatus, middleware.JWTMiddleware(), middleware.IsAdmin())
	users.Delete("/:id", handlers.DeleteUser, middleware.JWTMiddleware(), middleware.IsSelfOrAdmin())

	habits := api.Group("/habits")
	habits.Get("/", handlers.GetHabits, middleware.JWTMiddleware(), middleware.IsAdmin())
	habits.Get("/user/:id", handlers.GetUserHabits, middleware.JWTMiddleware(), middleware.IsSelfOrAdmin())
	habits.Get("/:id", handlers.GetHabit, middleware.JWTMiddleware())
	habits.Post("/", handlers.CreateHabit, middleware.JWTMiddleware())
	habits.Put("/:id", handlers.UpdateHabit, middleware.JWTMiddleware(), middleware.IsCreatorOrAdmin())
	habits.Delete("/:id", handlers.DeleteHabit, middleware.JWTMiddleware(), middleware.IsCreatorOrAdmin())

	notificationsSettings := api.Group("/notifications_settings")
	notificationsSettings.Get("/user/:id", handlers.GetUserNotificationSettings, middleware.JWTMiddleware(), middleware.IsSelfOrAdmin())
	notificationsSettings.Put("/user/:id", handlers.UpdateUserNotificationSettings, middleware.JWTMiddleware(), middleware.IsSelfOrAdmin())

	internal := app.Group("/internal")
	internalNotificationsSettings := internal.Group("/notifications_settings")
	internalNotificationsSettings.Get("/user/:id", handlers.GetUserNotificationSettings, middleware.APIKeyMiddleware())
}
