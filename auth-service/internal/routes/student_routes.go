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

	// 🔐 Protected base route
	r := app.Group("/api/v1/students", middleware.AuthMiddleware(cfg))

	// SINGLE student create  (only content_manager)
	r.Post("/", middleware.RequireRole("content_manager"), h.CreateStudent)

	// MULTIPLE students create (only content_manager)
	r.Post("/bulk", middleware.RequireRole("content_manager"), h.CreateManyStudents)

	// Anyone logged in can view students by batch
	r.Get("/batch/:batchID", h.GetByBatch)
}
