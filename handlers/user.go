package handlers

import (
	"context"
	"time"

	"auth-api/configs"
	"auth-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

func GetProfile(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	email := claims["email"].(string)

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := configs.DB.Collection("users").FindOne(ctx, bson.M{"email":email}).Decode(&user)
	if err != nil {
		return err
	}

	user.Password = "" // remove sensitive information

	return c.JSON(fiber.Map{
		"success": true,
		"user": user,
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	email := claims["email"].(string)

	var input struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"email": input.Email,
		},
	}

	result, err := configs.DB.Collection("users").UpdateOne(
		ctx,
		bson.M{"email": email},
		update,
	)

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile updated successfully",
	})
}

func DeleteProfile(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	email := claims["email"].(string)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := configs.DB.Collection("users").DeleteOne(ctx, bson.M{"email": email})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile deleted successfuly",
	})
}