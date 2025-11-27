package routes

import (
	"auth-service/internal/config"
	"auth-service/internal/handlers"
	"auth-service/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterBatchRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	h := handlers.NewBatchHandler(db)
	hy := handlers.NewBatchYearHandler(db)

	r := app.Group("/api/v1/batch-years", middleware.AuthMiddleware(cfg))

	r.Post("/", hy.CreateBatchYear)
	r.Get("/institution/:instID", hy.GetByInstitution)

	// Batches
	b := app.Group("/api/v1/batches", middleware.AuthMiddleware(cfg))
	b.Post("/", h.CreateBatch)
	b.Get("/year/:yearID", h.GetByBatchYear)
}
