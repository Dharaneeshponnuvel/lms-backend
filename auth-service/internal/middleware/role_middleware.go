package middleware

import (
	"auth-service/internal/auth"

	"github.com/gofiber/fiber/v2"
)

// Allow only specific role (ex: ADMIN)
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claimsAny := c.Locals("claims")
		if claimsAny == nil {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		claims, ok := claimsAny.(*auth.Claims)
		if !ok || claims.Role != role {
			return c.Status(403).JSON(fiber.Map{"error": "Forbidden: Admins only"})
		}

		return c.Next()
	}
}
