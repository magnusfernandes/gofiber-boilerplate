package userhandlers

import (
	"strings"

	"github.com/magnusfernandes/gofiber-boilerplate/database"
	"github.com/magnusfernandes/gofiber-boilerplate/helpers"
	"github.com/magnusfernandes/gofiber-boilerplate/models"
	"github.com/gofiber/fiber/v2"
)

func ListUsers(c *fiber.Ctx) error {
	user := database.FindUserAuth(c)
	if user.Role == "end_user" {
		return helpers.NotAuthorizedError(c, "Need administrator role to access!")
	}

	searchQuery := strings.Trim(c.Query("q"), " ")

	var users []models.User
	database.Database.
		Where("name ILIKE ? AND role != 'superadmin'", "%"+searchQuery+"%").
		Or("email ILIKE ? AND role != 'superadmin'", "%"+searchQuery+"%").
		Find(&users)
	var payload = make([]fiber.Map, 0)
	for _, user := range users {
		payload = append(payload, user.Json())
	}

	return c.JSON(helpers.BuildResponse(fiber.Map{
		"users": payload,
	}))
}
