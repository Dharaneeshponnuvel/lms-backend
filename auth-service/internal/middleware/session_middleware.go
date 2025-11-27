package middleware

import (
	"strings"
	"time"

	"auth-service/internal/auth"
	"auth-service/internal/config"
	"auth-service/internal/redis"

	"github.com/gofiber/fiber/v2"
)

func SessionMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := extractToken(c.Get("Authorization"))
		if token == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing token"})
		}

		claims, err := auth.VerifyToken(token, cfg)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}

		sessionKey := "session:" + claims.UserID
		storedToken, err := redis.RDB.Get(redis.Ctx, sessionKey).Result()
		if err != nil || storedToken != token {
			return c.Status(401).JSON(fiber.Map{"error": "session expired"})
		}

		// ❌ OLD: redis.RDB.Expire(redis.Ctx, sessionKey, cfg.AccessTokenTTL)
		// ✔ FIX:
		redis.RDB.Expire(redis.Ctx, sessionKey, time.Duration(cfg.AccessTokenTTL)*time.Second)

		c.Locals("claims", claims)
		return c.Next()
	}
}

func extractToken(header string) string {
	parts := strings.Split(header, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
