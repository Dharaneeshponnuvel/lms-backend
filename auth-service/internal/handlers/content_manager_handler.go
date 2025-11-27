package handlers

import (
	"auth-service/internal/auth"
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"auth-service/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContentManagerHandler struct {
	service *services.ContentManagerService
}

func NewContentManagerHandler(db *gorm.DB) *ContentManagerHandler {
	repo := repositories.NewContentManagerRepository(db)
	service := services.NewContentManagerService(repo)
	return &ContentManagerHandler{service: service}
}

// 🔥 POST /api/v1/content-managers/
func (h *ContentManagerHandler) CreateContentManager(c *fiber.Ctx) error {
	var req models.CreateContentManagerRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	claims := c.Locals("claims").(*auth.Claims)
	createdBy, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid creator ID"})
	}

	cmID, err := h.service.CreateWithUser(req, createdBy)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create content manager"})
	}

	return c.JSON(fiber.Map{
		"success":            true,
		"message":            "Content Manager created successfully",
		"content_manager_id": cmID,
	})
}

// 📌 POST /api/v1/content-managers/:id/batches
func (h *ContentManagerHandler) AssignBatch(c *fiber.Ctx) error {
	managerID := c.Params("id")

	var assign models.ContentManagerBatch
	if err := c.BodyParser(&assign); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid data"})
	}

	assign.ContentManagerID = uuid.MustParse(managerID)
	assign.AssignedBy = uuid.MustParse(c.Locals("claims").(*auth.Claims).UserID)

	if err := h.service.AssignBatch(&assign); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to assign batch"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Batch assigned successfully"})
}

// 📌 GET /api/v1/content-managers/:id
func (h *ContentManagerHandler) GetContentManagerByID(c *fiber.Ctx) error {
	userID := c.Params("id")

	manager, err := h.service.GetManagerDetails(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Content Manager not found"})
	}

	return c.JSON(fiber.Map{"success": true, "data": manager})
}

// 📌 GET /api/v1/content-managers/:id/batches
func (h *ContentManagerHandler) GetManagerBatches(c *fiber.Ctx) error {
	managerID := c.Params("id")

	batches, err := h.service.GetManagerBatches(managerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch batches"})
	}

	return c.JSON(fiber.Map{"success": true, "data": batches})
}
