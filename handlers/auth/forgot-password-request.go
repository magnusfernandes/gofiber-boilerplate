package authHandlers

import (
	log "github.com/sirupsen/logrus"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/magnusfernandes/gofiber-boilerplate/database"
	"github.com/magnusfernandes/gofiber-boilerplate/helpers"
	"github.com/magnusfernandes/gofiber-boilerplate/models"
)

func ForgotPasswordRequest(c *fiber.Ctx) error {
	jsonBody := struct {
		Phone string `json:"phone" validate:"required"`
	}{}

	//validation
	if err := c.BodyParser(&jsonBody); err != nil {
		return helpers.BadRequestError(c, "Missing attributes")
	}
	validate := validator.New()
	err := validate.Struct(jsonBody)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.StructNamespace() == "Phone" {
				return helpers.BadRequestError(c, "Please check phone number")
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

	client, err := helpers.InitPlivo()
	if err != nil {
		log.Error("Init plivo: ", err.Error())
		return helpers.BadRequestError(c, "Error initializing SMS gateway.")
	}

	code, err := helpers.GenerateCode(6)
	if err != nil {
		log.Error("Generate OTP: ", err.Error())
		return helpers.BadRequestError(c, "Error generating OTP.")
	}

	otpRequest := new(models.OTPRequest)
	otpRequest.Code = code
	otpRequest.UserId = user.Id
	transaction := database.Database.Begin()
	if err = transaction.Create(otpRequest).Error; err != nil {
		log.Info(err.Error())
		transaction.Rollback()
		return helpers.BadRequestError(c, "There was an error!")
	}

	_, err = helpers.SendSMS(client, user.Phone, "Your Let's Rate It verification code is "+code)
	if err != nil {
		log.Error("Send SMS: ", err.Error())
		transaction.Rollback()
		return helpers.BadRequestError(c, "Error sending SMS.")
	}
	transaction.Commit()

	return c.JSON(helpers.BuildResponse(fiber.Map{
		"requestId": otpRequest.Id,
	}))
}
