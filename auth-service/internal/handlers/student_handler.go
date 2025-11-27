package handlers

import (
	"auth-service/internal/auth"
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"auth-service/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type StudentHandler struct {
	service *services.StudentService
}

func NewStudentHandler(db *gorm.DB) *StudentHandler {
	repo := repositories.NewStudentRepository(db)
	service := services.NewStudentService(repo)
	return &StudentHandler{service}
}

func (h *StudentHandler) CreateStudent(c *fiber.Ctx) error {
	var st models.Student
	if err := c.BodyParser(&st); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	// ⛔ Check Batch ID is valid
	if st.BatchID == uuid.Nil {
		return c.Status(400).JSON(fiber.Map{"error": "BatchID is required"})
	}

	// 🔍 Check if Batch exists
	db := c.Locals("db").(*gorm.DB)
	var batch models.Batch
	if err := db.First(&batch, "id = ?", st.BatchID).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Batch ID"})
	}

	// 🔐 Get User ID from Token
	claims := c.Locals("claims").(*auth.Claims)
	userID, _ := uuid.Parse(claims.UserID)
	st.UserID = userID

	// 🧠 Set JSON properly
	if st.Demographics == nil {
		st.Demographics = datatypes.JSON([]byte("{}"))
	}

	if err := h.service.CreateStudent(&st); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Unable to create student", "details": err.Error()})
	}

	return c.Status(201).JSON(st)
}

func (h *StudentHandler) GetByBatch(c *fiber.Ctx) error {
	batchID := c.Params("batchID")
	students, err := h.service.GetByBatch(batchID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch students"})
	}
	return c.JSON(students)
}
