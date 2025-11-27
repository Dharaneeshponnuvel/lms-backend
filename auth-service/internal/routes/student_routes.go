package routes

import (
	"auth-service/internal/config"
	"auth-service/internal/handlers"
	"auth-service/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterStudentRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	h := handlers.NewStudentHandler(db)

	r := app.Group("/api/v1/students", middleware.AuthMiddleware(cfg))
	r.Post("/", h.CreateStudent)
	r.Get("/batch/:batchID", h.GetByBatch)
}
