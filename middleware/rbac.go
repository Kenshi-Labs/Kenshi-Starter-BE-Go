package middleware

import (
	"context"
	"time"

	"auth-api/configs"
	"auth-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

// RBACMiddleware checks if user has required permission
func RBACMiddleware(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// get user claims from context
		claims := c.Locals("user").(jwt.MapClaims)
		userRole := claims["role"].(string)

		// get role permissions from database
		var role models.Role
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := configs.DB.Collection("roles").FindOne(ctx, bson.M{"name":userRole}).Decode(&role)
		if err != nil {
			return fiber.NewError(fiber.StatusForbidden, "Role not found")
		}

		// check if user has the required permission`
		hasPermission := false
		for _, permission := range role.Permissions {
			if permission == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return fiber.NewError(fiber.StatusForbidden, "Insufficient permissions")
		}

		return c.Next()
	}
}