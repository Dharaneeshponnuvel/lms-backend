package middleware

import (
	"auth-service/internal/auth"

	"github.com/gofiber/fiber/v2"
)

// Require specific role (content_manager, admin,...)
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claimsAny := c.Locals("claims")
		if claimsAny == nil {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		claims, ok := claimsAny.(*auth.Claims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token data"})
		}

		// Compare role
		if claims.Role != role {
			return c.Status(403).JSON(fiber.Map{
				"error": "Forbidden: Only " + role + " can perform this action",
			})
		}

		return c.Next()
	}
}
