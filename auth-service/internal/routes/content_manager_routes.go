package routes

import (
	"auth-service/internal/config"
	"auth-service/internal/handlers"
	"auth-service/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterContentManagerRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {

	h := handlers.NewContentManagerHandler(db)

	// 🚀 SECURITY: ADMIN + INSTITUTION ONLY
	cmGroup := app.Group("/api/v1/content-managers",
		middleware.AuthMiddleware(cfg),
		middleware.SessionMiddleware(cfg),
		middleware.RequireRole("ADMIN"), // ✔ NOW VALID
		middleware.AuditMiddleware(cfg),
	)

	// 📌 ADMIN & INSTITUTION Routes (CREATE & ASSIGN)
	cmGroup.Post("/", h.CreateContentManager)
	cmGroup.Post("/:id/batches", h.AssignBatch)

	// ✔ ANY LOGGED IN USER CAN VIEW DETAILS
	app.Get("/api/v1/content-managers/:id",
		middleware.AuthMiddleware(cfg),
		middleware.SessionMiddleware(cfg),
		h.GetContentManagerByID,
	)

	app.Get("/api/v1/content-managers/:id/batches",
		middleware.AuthMiddleware(cfg),
		middleware.SessionMiddleware(cfg),
		h.GetManagerBatches,
	)
}
