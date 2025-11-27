package handlers

import (
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"auth-service/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RoleHandler struct {
	service *services.RoleService
}

func NewRoleHandler(db *gorm.DB) *RoleHandler {
	repo := repositories.NewRoleRepository(db)
	service := services.NewRoleService(repo)
	return &RoleHandler{service}
}

// 📌 GET /api/v1/roles
func (h *RoleHandler) GetAllRoles(c *fiber.Ctx) error {
	roles, err := h.service.GetRoles()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch roles"})
	}
	return c.JSON(roles)
}

// 📌 POST /api/v1/roles
func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var role models.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	if role.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Role name is required"})
	}

	err := h.service.CreateRole(&role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create role"})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "role": role})
}

// 📌 DELETE /api/v1/roles/:id
func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Role ID required"})
	}

	err := h.service.DeleteRole(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete role"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Role deleted"})
}
