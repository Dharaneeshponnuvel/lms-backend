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

type BatchYearHandler struct {
	service *services.BatchYearService
}

func NewBatchYearHandler(db *gorm.DB) *BatchYearHandler {
	repo := repositories.NewBatchYearRepository(db)
	service := services.NewBatchYearService(repo)
	return &BatchYearHandler{service}
}

func (h *BatchYearHandler) CreateBatchYear(c *fiber.Ctx) error {
	var by models.BatchYear
	if err := c.BodyParser(&by); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// 👇 Get from middleware (claims)
	claims := c.Locals("claims")
	if claims == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized - No Token"})
	}

	userClaims := claims.(*auth.Claims) // from JWT
	userID, _ := uuid.Parse(userClaims.UserID)

	by.CreatedBy = userID // ⭐ SET HERE
	if err := h.service.Create(&by); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create batch year", "details": err.Error()})
	}
	return c.Status(201).JSON(by)
}
func (h *BatchYearHandler) GetByInstitution(c *fiber.Ctx) error {
	instID := c.Params("instID") // ← MUST MATCH router
	years, err := h.service.GetInstitutionBatchYears(instID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "No batch years found"})
	}
	return c.JSON(years)
}
