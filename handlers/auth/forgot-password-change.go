package authHandlers

import (
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/magnusfernandes/gofiber-boilerplate/database"
	"github.com/magnusfernandes/gofiber-boilerplate/helpers"
	"github.com/magnusfernandes/gofiber-boilerplate/models"
)

func ForgotPasswordChange(c *fiber.Ctx) error {

	jsonBody := struct {
		Password  string `json:"password" validate:"required"`
		RequestId string `json:"requestId" validate:"required"`
		Code      string `json:"otp" validate:"required"`
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

	var otpRequest models.OTPRequest
	err = database.Database.
		First(&otpRequest, "id = ? AND code = ?", jsonBody.RequestId, jsonBody.Code).Error
	if err != nil {
		log.Error(err)
		return helpers.BadRequestError(c, "There was an error!")
	}

	if otpRequest.IsExpired || time.Since(otpRequest.CreatedAt).Minutes() > 3 {
		return helpers.BadRequestError(c, "OTP request expired!")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(jsonBody.Password), 8)
	if err != nil {
		log.Info(err.Error())
		return helpers.BadRequestError(c, "There was an error!")
	}

	transaction := database.Database.Begin()
	err = transaction.Model(&models.User{}).
		Where("id = ?", otpRequest.UserId).
		Update("password_hash", string(hashedPassword)).Error
	if err != nil {
		log.Error(err)
		transaction.Rollback()
		return helpers.BadRequestError(c, "There was an error!")
	}

	err = transaction.Model(&models.OTPRequest{}).
		Where("id = ?", otpRequest.Id).
		Update("is_expired", true).Error
	if err != nil {
		log.Error(err)
		transaction.Rollback()
		return helpers.BadRequestError(c, "There was an error!")
	}

	transaction.Commit()

	return c.JSON(helpers.BuildResponse("Password changed successfully!"))
}
