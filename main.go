package main

import (
	"log"

	"auth-api/configs"
	"auth-api/handlers"
	"auth-api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	config := configs.AppConfig

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// initialize mongodb connection
	configs.ConnectDB()

	// setup routes
	api := app.Group("/api")

	// auth routes with validation
	auth := api.Group("/auth")
	auth.Post("/signup", middleware.ValidateUserInput(),handlers.SignUp)
	auth.Post("/signin", middleware.ValidateUserInput(),handlers.SignIn)
	auth.Post("/refresh", handlers.RefreshToken)

	// Password reset routes
	auth.Post("/forgot-password", handlers.RequestPasswordReset)
	auth.Post("/reset-password", handlers.ResetPassword)

	// protected routes
	// user routes with authentication and RBAC
	user := api.Group("/user", middleware.AuthMiddleware())
	user.Get("/profile", middleware.RBACMiddleware("read:profile"), handlers.GetProfile)
	user.Put("/profile", middleware.RBACMiddleware("update:profile"), handlers.UpdateProfile)
	user.Delete("profile", middleware.RBACMiddleware("delete:profile"), handlers.DeleteProfile)


	port := config.APIPort

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}