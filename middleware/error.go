package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// helper function to check for MongoDB duplicate key error
func isDuplicateKeyError(err error) bool {
	var mongoErr mongo.WriteException
	if errors.As(err, &mongoErr) {
		for _, writeErr := range mongoErr.WriteErrors{
			if writeErr.Code == 11000 {
				return true
			}
		}
	}
	return false
}

// helper function to check for MongoDB timeout error
func isTimeoutError(err error) bool {
	return	strings.Contains(err.Error(), "deadline exceeded") || 
			strings.Contains(err.Error(), "context deadline exceeded")
}

// ErrorHandler handles all errors in the application
func ErrorHandler(c *fiber.Ctx, err error) error {
	// default error
	code := fiber.StatusInternalServerError
	message := "Internal server error"

	// check error type and set appropriate status code and message
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		code = fiber.StatusNotFound
		message = "Resource not found"

	//handle mongodb duplicate key error
	case isDuplicateKeyError(err):
		code = fiber.StatusConflict
		message = "Resource already exists"
	
	case isTimeoutError(err):
		code = fiber.StatusRequestTimeout
		message="Request timed out"
	
	case err.Error()=="invalid token":
		code = fiber.StatusUnauthorized
		message = "Invalid or expired token"
	default:
		// check if it's a fiber error
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = e.Message
		}
	}

	// return JSON response
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error": message,
		"status": code,
	})
}

