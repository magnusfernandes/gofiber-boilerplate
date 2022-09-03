package helpers

import "github.com/gofiber/fiber/v2"

func BuildResponse(v interface{}) fiber.Map {
	res := fiber.Map{
		"type":    "success",
		"message": v,
	}
	return res
}

func BadRequestError(c *fiber.Ctx, v interface{}) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"type":    "error",
		"message": v,
	})
}

func ResourceNotFoundError(c *fiber.Ctx, resource string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"type":    "error",
		"message": resource + " not found",
	})
}

func NotAuthorizedError(c *fiber.Ctx, v interface{}) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"type":    "error",
		"message": v,
	})
}

func CustomError(c *fiber.Ctx, code int, v interface{}) error {
	return c.Status(code).JSON(fiber.Map{
		"type":    "error",
		"message": v,
	})
}
