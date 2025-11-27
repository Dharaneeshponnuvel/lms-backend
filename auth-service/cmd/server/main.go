package main

import (
	"log"

	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/redis"
	"auth-service/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	redis.Connect(cfg)

	app := fiber.New()

	// CORS Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173", // React frontend
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	// Attach DB to context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	// ⭐ FIX HERE — Pass DB also!
	routes.RegisterAuthRoutes(app, db, cfg)
	routes.RegisterRoleRoutes(app, db)
	routes.RegisterStudentRoutes(app, db, cfg)
	routes.RegisterInstitutionRoutes(app, db, cfg)
	routes.RegisterSessionRoutes(app, db)
	routes.RegisterBatchRoutes(app, db, cfg)
	routes.RegisterContentManagerRoutes(app, db, cfg)

	log.Println("Auth service running on port", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
