package userhandlers

import (
	"github.com/magnusfernandes/gofiber-boilerplate/database"
	"github.com/magnusfernandes/gofiber-boilerplate/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUserDetails(c *fiber.Ctx) error {
	userId, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return helpers.BadRequestError(c, "Invalid UUID!")
	}

	user := database.FindUserById(userId)

	return c.JSON(helpers.BuildResponse(fiber.Map{
		"user": user.Json(),
	}))
}
