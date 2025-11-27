package routes

import (
	"auth-service/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoleRoutes(app *fiber.App, db *gorm.DB) {
	h := handlers.NewRoleHandler(db)

	r := app.Group("/api/v1/roles")
	r.Get("/", h.GetAllRoles)
	r.Post("/", h.CreateRole)
	r.Delete("/:id", h.DeleteRole)
}
