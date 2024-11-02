package utils

import (
    "crypto/rand"
    "encoding/base64"
    "time"
    "context"
    
    "auth-api/configs"
    "auth-api/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func GenerateRefreshToken(userID primitive.ObjectID) (*models.RefreshToken, error) {
    // Generate random token
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return nil, err
    }
    
    refreshToken := &models.RefreshToken{
        Token:     base64.URLEncoding.EncodeToString(b),
        UserID:    userID,
        ExpiresAt: time.Now().Add(time.Hour * 24 * 7), // 7 days
        CreatedAt: time.Now(),
    }
    
    // Store in database
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    _, err := configs.DB.Collection("refresh_tokens").InsertOne(ctx, refreshToken)
    if err != nil {
        return nil, err
    }
    
    return refreshToken, nil
}

func ValidateRefreshToken(token string) (*models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    var refreshToken models.RefreshToken
    err := configs.DB.Collection("refresh_tokens").FindOne(ctx, bson.M{
        "token": token,
        "expires_at": bson.M{"$gt": time.Now()},
    }).Decode(&refreshToken)
    
    if err != nil {
        return nil, err
    }
    
    // Get user
    var user models.User
    err = configs.DB.Collection("users").FindOne(ctx, bson.M{
        "_id": refreshToken.UserID,
    }).Decode(&user)
    
    return &user, err
}