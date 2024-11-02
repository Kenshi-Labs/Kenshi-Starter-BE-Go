package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	// "go.mongodb.org/mongo-driver/bson/primitive"
	"auth-api/configs"
	"auth-api/models"
	"auth-api/utils"
)

func RequestPasswordReset(c *fiber.Ctx) error {
    var input struct {
        Email string `json:"email"`
    }
    
    if err := c.BodyParser(&input); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
    }
    
    // Generate reset token
    b := make([]byte, 16)
    if _, err := rand.Read(b); err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Could not generate reset token")
    }
    token := hex.EncodeToString(b)
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Find user
    var user models.User
    err := configs.DB.Collection("users").FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
    if err != nil {
        return fiber.NewError(fiber.StatusNotFound, "User not found")
    }
    
    // Create password reset record
    resetRecord := models.PasswordReset{
        UserID:    user.ID,
        Token:     token,
        ExpiresAt: time.Now().Add(time.Hour), // 1 hour expiry
        CreatedAt: time.Now(),
    }
    
    _, err = configs.DB.Collection("password_resets").InsertOne(ctx, resetRecord)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Could not create reset token")
    }
    
    // Send email
    resetLink := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("APP_URL"), token)
    emailBody := fmt.Sprintf("Click the following link to reset your password: %s", resetLink)
    
    err = utils.SendEmail(user.Email, "Password Reset", emailBody)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Could not send reset email")
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "message": "Password reset email sent",
    })
}

func ResetPassword(c *fiber.Ctx) error {
    var input struct {
        Token    string `json:"token"`
        Password string `json:"password"`
    }
    
    if err := c.BodyParser(&input); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Find reset token
    var resetRecord models.PasswordReset
    err := configs.DB.Collection("password_resets").FindOne(ctx, bson.M{
        "token": input.Token,
        "expires_at": bson.M{"$gt": time.Now()},
    }).Decode(&resetRecord)
    
    if err != nil {
        return fiber.NewError(fiber.StatusNotFound, "Invalid or expired reset token")
    }
    
    // Hash new password
    hashedPassword, err := utils.HashPassword(input.Password)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Could not hash password")
    }
    
    // Update user password
    _, err = configs.DB.Collection("users").UpdateOne(
        ctx,
        bson.M{"_id": resetRecord.UserID},
        bson.M{"$set": bson.M{"password": hashedPassword}},
    )
    
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Could not update password")
    }
    
    // Delete used reset token
    _, err = configs.DB.Collection("password_resets").DeleteOne(ctx, bson.M{"token": input.Token})
    if err != nil {
        // Log error but don't return it to user
        log.Printf("Could not delete used reset token: %v", err)
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "message": "Password reset successfully",
    })
}