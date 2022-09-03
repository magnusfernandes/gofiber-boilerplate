package routes

import (
	"github.com/gofiber/fiber/v2"
	authHandlers "github.com/magnusfernandes/gofiber-boilerplate/handlers/auth"
	profileHandlers "github.com/magnusfernandes/gofiber-boilerplate/handlers/profile"
	"github.com/magnusfernandes/gofiber-boilerplate/middlewares"
)

func AuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/login", authHandlers.LoginUser)
	auth.Post("/signup", authHandlers.SignupUser)
	auth.Get("/validate", middlewares.Protected(), authHandlers.ValidateUser)
	auth.Patch("/profile", middlewares.Protected(), profileHandlers.UpdateProfile)
	auth.Patch("/change-password", middlewares.Protected(), profileHandlers.ChangePassword)
	auth.Post("/request-forgot-password", authHandlers.ForgotPasswordRequest)
	auth.Post("/change-forgot-password", authHandlers.ForgotPasswordChange)
}
