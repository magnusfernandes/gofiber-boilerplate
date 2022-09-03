package userhandlers

import (
	log "github.com/sirupsen/logrus"

	"github.com/magnusfernandes/gofiber-boilerplate/database"
	"github.com/magnusfernandes/gofiber-boilerplate/helpers"
	"github.com/magnusfernandes/gofiber-boilerplate/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func EditUser(c *fiber.Ctx) error {
	userId, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return helpers.BadRequestError(c, "Invalid UUID!")
	}

	jsonBody := struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
		Phone string `json:"phone" validate:"required"`
	}{}
	if err := c.BodyParser(&jsonBody); err != nil {
		log.Info(err)
		return helpers.BadRequestError(c, "Error parsing body!")
	}

	validate := validator.New()
	err = validate.Struct(jsonBody)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.StructNamespace() == "Email" && err.Tag() == "email" {
				return helpers.BadRequestError(c, "Please check the email")
			}
		}
		log.Info(err)
		return helpers.BadRequestError(c, "Please check your request!")
	}

	var count int64
	database.Database.Model(&models.User{}).
		Where("email = ? AND id != ?", jsonBody.Email, userId).
		Or("phone = ? AND id != ?", jsonBody.Phone, userId).
		Count(&count)
	if count > 0 {
		return helpers.BadRequestError(c, "Email or Phone already exists!")
	}

	var user models.User
	database.Database.
		Model(&models.User{}).Where("id = ?", userId).
		Updates(models.User{Name: jsonBody.Name, Email: jsonBody.Email, Phone: jsonBody.Phone})

	user = database.FindUserById(userId)
	return c.JSON(helpers.BuildResponse(fiber.Map{
		"user": user.Json(),
	}))
}
