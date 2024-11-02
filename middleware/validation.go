package middleware

import (
	"auth-api/utils"

	"github.com/gofiber/fiber/v2"
)

// validateuserinput validates request body for auth routes
func ValidateUserInput() fiber.Handler{
	return func(c *fiber.Ctx) error {
		switch c.Path() {
		case "/api/auth/signup", "/api/auth/signin":
			var input struct {
				Email string `json:"email"`
				Password string `json:"password"`
			}
			
			if err := c.BodyParser(&input); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
			}

			if !utils.IsValidEmail(input.Email) {
				return fiber.NewError(fiber.StatusBadRequest, "Invalid email format")
			}

			if len(input.Password) < 8 {
				return fiber.NewError(fiber.StatusBadRequest, "Password must be at least 8 characters long")
			}
		}
		return c.Next()
	}
}