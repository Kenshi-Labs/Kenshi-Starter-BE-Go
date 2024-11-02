package middleware

import (
	"strings"

	"auth-api/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware checks for valid JWT token
func AuthMiddleware() fiber.Handler{
	return func(c *fiber.Ctx)error {
		// get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
		}

		// Check if it's a bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format")
		}

		// validate token
		claims, err := utils.ValidateJWT(tokenParts[1])
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}

		// store user claims in context
		c.Locals("user", claims)
		return c.Next()
	}
}