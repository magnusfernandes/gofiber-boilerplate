package userhandlers

import (
	"github.com/magnusfernandes/gofiber-boilerplate/database"
	"github.com/magnusfernandes/gofiber-boilerplate/helpers"
	"github.com/magnusfernandes/gofiber-boilerplate/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const CommonPassword = "changeme"

func CreateNewUser(c *fiber.Ctx) error {
	jsonBody := struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
		Phone string `json:"phone" validate:"required"`
		Role  string `json:"role" validate:"required"`
	}{}
	if err := c.BodyParser(&jsonBody); err != nil {
		log.Info(err)
		return helpers.BadRequestError(c, "Error parsing body!")
	}

	validate := validator.New()
	err := validate.Struct(jsonBody)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.StructNamespace() == "Email" && err.Tag() == "email" {
				return helpers.BadRequestError(c, "Please check the email")
			}
		}
		return helpers.BadRequestError(c, "Please check your request!")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(CommonPassword), 8)
	if err != nil {
		log.Info(err.Error())
		return helpers.BadRequestError(c, "There was an error!")
	}

	user := new(models.User)

	user.Name = jsonBody.Name
	user.Email = jsonBody.Email
	user.Phone = jsonBody.Phone
	user.PasswordHash = string(hashedPassword)
	user.Role = jsonBody.Role

	var count int64
	database.Database.Model(&models.User{}).Where("email = ?", jsonBody.Email).Or("phone = ?", jsonBody.Phone).Count(&count)

	if count > 0 {
		return helpers.BadRequestError(c, "Email or Phone already exists!")
	}

	if err = database.Database.Create(user).Error; err != nil {
		log.Info(err.Error())
	}
	return c.JSON(helpers.BuildResponse(fiber.Map{
		"user": user.Json(),
	}))
}
