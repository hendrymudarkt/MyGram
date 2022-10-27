package middlewares

import (
	"MyGram/helpers"

	"github.com/gofiber/fiber/v2"
)

func Authentication() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		verifyToken, err := helpers.VerifyToken(c)
		_ = verifyToken

		if err != nil {
			c.Status(500).JSON(fiber.Map{
				"error": "Unauthenticated",
				"message": err.Error(),
			})
			return
		}
		c.Locals("userData", verifyToken)
		c.Next()

		return
	}
}