package profileHandlers

import (
	log "github.com/sirupsen/logrus"

	"github.com/magnusfernandes/gofiber-boilerplate/database"
	"github.com/magnusfernandes/gofiber-boilerplate/helpers"
	"github.com/magnusfernandes/gofiber-boilerplate/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func ChangePassword(c *fiber.Ctx) error {
	user := database.FindUserAuth(c)

	jsonBody := struct {
		Password    string `json:"password" validate:"required"`
		OldPassword string `json:"oldPassword" validate:"required"`
	}{}
	if err := c.BodyParser(&jsonBody); err != nil {
		log.Info(err)
		return helpers.BadRequestError(c, "Error parsing body!")
	}

	validate := validator.New()
	err := validate.Struct(jsonBody)
	if err != nil {
		return helpers.BadRequestError(c, "Please check your request!")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(jsonBody.OldPassword))
	if err != nil {
		log.Info(err)
		return helpers.BadRequestError(c, "Wrong password entered!")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(jsonBody.Password), 8)
	if err != nil {
		log.Info(err.Error())
		return helpers.BadRequestError(c, "There was an error!")
	}

	err = database.Database.Model(&models.User{}).
		Where("id = ?", user.Id).
		Update("password_hash", string(hashedPassword)).Error
	if err != nil {
		log.Error(err)
	}

	return c.JSON(helpers.BuildResponse(fiber.Map{
		"user": user.Json(),
	}))
}
