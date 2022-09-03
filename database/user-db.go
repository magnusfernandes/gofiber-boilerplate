package database

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/magnusfernandes/gofiber-boilerplate/models"
)

func ListUsers() []models.User {
	var users []models.User
	Database.Find(&users)
	return users
}

func FindUserById(id uuid.UUID) models.User {
	var user models.User
	Database.Preload("Organisations").First(&user, "id = ?", id)
	return user
}

func FindUserAuth(ctx *fiber.Ctx) models.User {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	idString := claims["id"].(string)
	id, _ := uuid.Parse(idString)
	return FindUserById(id)
}

func FindUserByPhone(phone string) (models.User, error) {
	var user models.User
	Database.Where("phone=?", phone).First(&user)
	if len(user.Id.String()) > 0 {
		return user, nil
	}
	return user, errors.New("no user found")
}

func FindUserByEmail(email string) (models.User, error) {
	var user models.User
	Database.Where("email=?", email).First(&user)
	if len(user.Id.String()) > 0 {
		return user, nil
	}
	return user, errors.New("no user found")
}

func DeleteUser(id uuid.UUID) error {
	var user models.User
	Database.First(&user, "id = ?", id)
	return Database.Delete(&user).Error
}
