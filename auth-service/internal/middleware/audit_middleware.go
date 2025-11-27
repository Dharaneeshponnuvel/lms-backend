package middleware

import (
	"encoding/json"

	"auth-service/internal/auth"
	"auth-service/internal/config"
	"auth-service/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AuditMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		db, ok := c.Locals("db").(*gorm.DB)
		if !ok {
			return err
		}

		// 🛡 SAFE — Check claims properly (avoid panic)
		var userID *uuid.UUID = nil
		if cl := c.Locals("claims"); cl != nil {
			if claims, ok := cl.(*auth.Claims); ok {
				if id, err := uuid.Parse(claims.UserID); err == nil {
					userID = &id
				}
			}
		}

		metaMap := map[string]interface{}{
			"method": c.Method(),
			"status": c.Response().StatusCode(),
		}
		meta, _ := json.Marshal(metaMap)

		logEntry := models.AuditLog{
			UserID:    userID,
			Action:    c.Method() + " " + c.Path(),
			Route:     c.Path(),
			IP:        c.IP(),
			UserAgent: string(c.Request().Header.UserAgent()),
			Metadata:  meta,
		}

		go db.Create(&logEntry) // 🔥 async insert

		return err
	}
}
