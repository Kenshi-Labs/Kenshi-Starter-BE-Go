package handlers

import (
	"context"
	"log"
	"time"

	"auth-api/configs"
	"auth-api/models"
	"auth-api/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUp(c *fiber.Ctx) error {
	var input struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Could not hash password")
	}

	user := models.User{
		Email: input.Email,
		Password : hashedPassword,
		Role: "user",// default role

	}

	// insert user
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := configs.DB.Collection("users").InsertOne(ctx, user)
	if err != nil {
		// check for duplicate email
		if mongo.IsDuplicateKeyError(err) {
			return fiber.NewError(fiber.StatusConflict, "Email already exists")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Could not create user")
	}

	user.Password = "" // remove password from response

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"user": user,
		"id": result.InsertedID,
	})
}

func SignIn(c *fiber.Ctx) error {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    
    if err := c.BodyParser(&input); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
    }
    
    // Find user
    var user models.User
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    err := configs.DB.Collection("users").FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
    if err != nil {
        return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
    }
    
    // Verify password
    if !utils.CheckPassword(input.Password, user.Password) {
        return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
    }
    
    // Generate access token
    accessToken, err := utils.GenerateJWT(user)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Could not generate access token")
    }
    
    // Generate refresh token
    refreshToken, err := utils.GenerateRefreshToken(user.ID)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Could not generate refresh token")
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "access_token": accessToken,
        "refresh_token": refreshToken.Token,
    })
}

func RefreshToken(c *fiber.Ctx) error {
    var input struct {
        RefreshToken string `json:"refresh_token"`
    }
    
    if err := c.BodyParser(&input); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
    }
    
    // Validate refresh token
    user, err := utils.ValidateRefreshToken(input.RefreshToken)
    if err != nil {
        return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
    }
    
    // Generate new access token
    accessToken, err := utils.GenerateJWT(*user)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Could not generate access token")
    }
    
    // Generate new refresh token
    newRefreshToken, err := utils.GenerateRefreshToken(user.ID)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Could not generate refresh token")
    }
    
    // Delete old refresh token
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    _, err = configs.DB.Collection("refresh_tokens").DeleteOne(ctx, bson.M{
        "token": input.RefreshToken,
    })
    if err != nil {
        log.Printf("Error deleting old refresh token: %v", err)
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "access_token": accessToken,
        "refresh_token": newRefreshToken.Token,
    })
}