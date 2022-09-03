package routes

import (
	"github.com/gofiber/fiber/v2"
	userHandlers "github.com/magnusfernandes/gofiber-boilerplate/handlers/users"
	"github.com/magnusfernandes/gofiber-boilerplate/middlewares"
)

func UserRoutes(app *fiber.App) {
	app.Post("/users/list", middlewares.Protected(), userHandlers.ListUsers)
	app.Post("/users", middlewares.Protected(), userHandlers.CreateNewUser)
	app.Get("/users/:userId", middlewares.Protected(), userHandlers.GetUserDetails)
	app.Patch("/users/:userId", middlewares.Protected(), userHandlers.EditUser)
	app.Delete("/users/:userId", middlewares.Protected(), userHandlers.DeleteUser)
}
