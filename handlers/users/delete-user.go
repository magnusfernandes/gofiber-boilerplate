package userhandlers

import (
	"github.com/magnusfernandes/gofiber-boilerplate/database"
	"github.com/magnusfernandes/gofiber-boilerplate/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func DeleteUser(c *fiber.Ctx) error {
	user := database.FindUserAuth(c)

	if user.Role == "end_user" {
		return helpers.NotAuthorizedError(c, "Need administrator role to access!")
	}

	userId, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return helpers.BadRequestError(c, "Invalid UUID!")
	}

	if err := database.DeleteUser(userId); err != nil {
		log.Info(err)
		return helpers.BadRequestError(c, "There was an error!")
	}

	return c.JSON(helpers.BuildResponse("User deleted successfully!"))
}
