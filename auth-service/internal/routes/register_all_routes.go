// internal/routes/register_all_routes.go
package routes

import (
	"auth-service/internal/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterAllRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	RegisterAuthRoutes(app, db, cfg)        // middleware + JWT
	RegisterRoleRoutes(app, db)             // roles
	RegisterStudentRoutes(app, db, cfg)     // students
	RegisterSessionRoutes(app, db)          // sessions
	RegisterInstitutionRoutes(app, db, cfg) // institutions
	RegisterBatchRoutes(app, db, cfg)       // batch + batch year
}
