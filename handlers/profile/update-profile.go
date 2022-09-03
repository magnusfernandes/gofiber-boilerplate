package profileHandlers

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/magnusfernandes/gofiber-boilerplate/database"
	"github.com/magnusfernandes/gofiber-boilerplate/helpers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func UpdateProfile(c *fiber.Ctx) error {
	user := database.FindUserAuth(c)

	jsonBody := struct {
		Name      string    `json:"name" validate:"required"`
		Email     string    `json:"email" validate:"required,email"`
		Gender    string    `json:"gender"`
		BirthDate time.Time `json:"birthDate"`
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

	user.Name = jsonBody.Name
	user.Email = jsonBody.Email

	if jsonBody.BirthDate.Year() > 1930 {
		user.BirthDate.Valid = true
		user.BirthDate.Time = jsonBody.BirthDate
	}

	if len(jsonBody.Gender) > 0 {
		user.Gender = jsonBody.Gender
	}

	database.Database.Save(&user)

	return c.JSON(helpers.BuildResponse(fiber.Map{
		"user": user.Json(),
	}))
}
