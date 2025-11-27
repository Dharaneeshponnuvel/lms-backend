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

type BatchHandler struct {
	service *services.BatchService
}

func NewBatchHandler(db *gorm.DB) *BatchHandler {
	repo := repositories.NewBatchRepository(db) // 1️Create Repo
	service := services.NewBatchService(repo)   // 2️ Pass Repo to Service
	return &BatchHandler{service: service}      // 3️ Inject Correctly
}

func (h *BatchHandler) CreateBatch(c *fiber.Ctx) error {
	var batch models.Batch
	if err := c.BodyParser(&batch); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// 🔐 VALIDATE JWT → GET USER ID
	claims := c.Locals("claims")
	if claims == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userClaims := claims.(*auth.Claims)

	createdBy, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	batch.CreatedBy = createdBy // <-- FIXED!

	err = h.service.Create(&batch)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create batch", "details": err.Error()})
	}

	return c.Status(201).JSON(batch)
}

func (h *BatchHandler) GetByBatchYear(c *fiber.Ctx) error {
	batchYearID := c.Params("yearID") // ← MATCH with router
	batches, err := h.service.GetByBatchYear(batchYearID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "No batches found"})
	}
	return c.JSON(batches)
}
