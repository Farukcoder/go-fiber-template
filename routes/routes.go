package routes

import (
	"go-fiber-template/controllers"
	"go-fiber-template/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Auth routes
	authController := controllers.NewAuthController(db)

	auth := app.Group("/api/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)

	// Protected routes example
	api := app.Group("/api")
	api.Use(middleware.Protected())

	// Add protected routes here
	api.Get("/profile", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Protected route",
			"user":    c.Locals("user"),
		})
	})
}
