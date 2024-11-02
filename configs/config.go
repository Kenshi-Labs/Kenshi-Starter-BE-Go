package configs

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoURI string
	JWTSecret string
	APIPort string
}

var (
	DB *mongo.Database
	AppConfig *Config
)

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// try to load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		MongoURI: getEnv("MONGODB_URI","mongodb://localhost:27017/authdb"),
		JWTSecret: getEnv("JWT_SECRET",""),
		APIPort: getEnv("API_PORT","3000"),
	}

	// validate required configurations
	if config.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return config
}

// getEnv retrieves an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func ConnectDB() {
	// set client options
	clientOptions := options.Client().ApplyURI(AppConfig.MongoURI)

	// Connect to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:",err)
	}

	// Ping the database
	err = client.Ping(ctx,nil)
	if err != nil {
		log.Fatal("Error Pinging MongoDB:", err)
	}

	log.Println("Connected to MongoDB!")

	// set database
	DB = client.Database("auth_db")
}

func init() {
	AppConfig = LoadConfig()
}