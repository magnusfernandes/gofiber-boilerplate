package authHandlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/magnusfernandes/gofiber-boilerplate/database"
	"github.com/magnusfernandes/gofiber-boilerplate/helpers"
)

func ValidateUser(c *fiber.Ctx) error {
	user := database.FindUserAuth(c)
	if len(user.Name) == 0 {
		return helpers.ResourceNotFoundError(c, "User not found!")
	}
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return helpers.NotAuthorizedError(c, "Invalid token!")
	}

	payload := fiber.Map{
		"user": user.Json(),
	}
	return c.JSON(helpers.BuildResponse(payload))
}
