package routes

import (
	"auth-service/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterSessionRoutes(app *fiber.App, db *gorm.DB) {
	h := handlers.NewSessionHandler(db)

	r := app.Group("/api/v1/sessions")
	r.Post("/", h.CreateSession)
	r.Post("/invalidate/:userID", h.InvalidateOldSessions)
	r.Get("/active/:userID", h.GetActiveSessions)
	r.Post("/logout/:userID/:token", h.LogoutSession)
}
