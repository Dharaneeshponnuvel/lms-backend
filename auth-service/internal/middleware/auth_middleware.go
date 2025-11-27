package middleware

import (
	"strings"

	"auth-service/internal/auth"
	"auth-service/internal/config"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing token"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token format"})
		}

		tokenStr := parts[1]

		// ❌ Old: auth.ParseToken(cfg, tokenStr)
		// ✅ FIX:
		claims, err := auth.VerifyToken(tokenStr, cfg)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals("claims", claims)
		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}
