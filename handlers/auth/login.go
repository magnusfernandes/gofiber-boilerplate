package authHandlers

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/magnusfernandes/gofiber-boilerplate/database"
	"github.com/magnusfernandes/gofiber-boilerplate/helpers"
	"github.com/magnusfernandes/gofiber-boilerplate/middlewares"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(c *fiber.Ctx) error {
	jsonBody := struct {
		Phone    string `json:"phone" validate:"required"`
		Password string `json:"password" validate:"required,min=6"`
	}{}

	//validation
	if err := c.BodyParser(&jsonBody); err != nil {
		return helpers.BadRequestError(c, "Missing attributes")
	}
	validate := validator.New()
	err := validate.Struct(jsonBody)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.StructNamespace() == "Password" && err.Tag() == "min" {
				return helpers.BadRequestError(c, "Password should be minimum 6 characters long")
			}
		}
	}

	user, err := database.FindUserByPhone(jsonBody.Phone)
	if err != nil {
		log.Info(err)
		return helpers.BadRequestError(c, "There was an error!")
	}

	if len(user.Name) == 0 {
		return helpers.NotAuthorizedError(c, "Username/password is incorrect!")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(jsonBody.Password))
	if err != nil {
		log.Info(err)
		return helpers.BadRequestError(c, "Wrong password entered!")
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.Id
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(middlewares.SECRET))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	payload := fiber.Map{
		"user":  user.Json(),
		"token": t,
	}
	return c.JSON(helpers.BuildResponse(payload))
}
