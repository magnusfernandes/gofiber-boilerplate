package middlewares

import (
	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	jwtWare "github.com/gofiber/jwt/v3"
	"github.com/magnusfernandes/gofiber-boilerplate/helpers"
)

const SECRET = "1ZVtHIAjEPcNrSHQ3lyWdqxFlsSA81YhdbXK7D57c0L19kxq27CqIQvkzxkPcMj6YARFVuvE8PLOeTGZbbjoWU904vdSCOA2dwaMd05cmnWt3iS4Xk0w9z1kFLtzcF"

func Protected() func(*fiber.Ctx) error {
	return jwtWare.New(jwtWare.Config{
		SigningKey:   []byte(SECRET),
		ErrorHandler: jwtError,
		AuthScheme:   "JWT",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	log.Info("jwtError:", err)
	if err.Error() == "Missing or malformed JWT" {
		return helpers.NotAuthorizedError(c, "Missing or malformed JWT")
	} else {
		return helpers.NotAuthorizedError(c, "Invalid or expired JWT")
	}
}
