package routes

import (
	"auth-service/internal/config"
	"auth-service/internal/handlers"
	"auth-service/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterInstitutionRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	h := handlers.NewInstitutionHandler(db)

	r := app.Group("/api/v1/institutions",
		middleware.AuthMiddleware(cfg),    // ✔ JWT verified
		middleware.SessionMiddleware(cfg), // ✔ Verify Redis Session
		middleware.RequireRole("ADMIN"),   // 🚫 Only ADMIN allowed
		middleware.AuditMiddleware(cfg),   // 📌 Audit Logs
	)

	r.Post("/", h.CreateInstitution) // ADMIN ONLY
	r.Get("/:id", h.GetByID)         // Anyone logged-in can view
}
