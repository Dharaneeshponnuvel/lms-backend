package routes

import (
	"auth-service/internal/config"
	"auth-service/internal/handlers"
	"auth-service/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	authGroup := app.Group("/api/v1/auth")

	authHandler := handlers.NewAuthHandler(cfg)
	roleHandler := handlers.NewRoleHandler(db)

	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/logout", middleware.AuthMiddleware(cfg), authHandler.Logout)
	authGroup.Get("/verify", middleware.AuthMiddleware(cfg), authHandler.Verify)
	authGroup.Get("/roles", roleHandler.GetAllRoles)
	authGroup.Post("/refresh", authHandler.RefreshToken)

	// Global Logging
	app.Use(middleware.AuditMiddleware(cfg))
}
