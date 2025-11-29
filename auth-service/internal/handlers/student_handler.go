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

type StudentHandler struct {
	service *services.StudentService
}

func NewStudentHandler(db *gorm.DB) *StudentHandler {
	repo := repositories.NewStudentRepository(db)
	service := services.NewStudentService(repo)
	return &StudentHandler{service}
}

// ---------- REQUEST STRUCTS ----------
type StudentInput struct {
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Name       string    `json:"name"`
	RollNumber string    `json:"roll_number"`
	Branch     string    `json:"branch"`
	BatchID    uuid.UUID `json:"batch_id"`
}

type CreateStudentBatchInput struct {
	Students []StudentInput `json:"students"`
}

// ==============================================
//
//	CREATE **ONE** STUDENT   (POST /api/v1/students)
//
// ==============================================
func (h *StudentHandler) CreateStudent(c *fiber.Ctx) error {
	var st StudentInput

	if err := c.BodyParser(&st); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	db := c.Locals("db").(*gorm.DB)
	claims := c.Locals("claims").(*auth.Claims)

	if claims.Role != "institution" {
		return c.Status(403).JSON(fiber.Map{"error": "Only institutions can create students"})
	}

	if st.BatchID == uuid.Nil {
		return c.Status(400).JSON(fiber.Map{"error": "BatchID is required"})
	}

	// validate batch
	var batch models.Batch
	if err := db.First(&batch, "id = ?", st.BatchID).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Batch ID"})
	}

	// get role "student"
	var studentRole models.Role
	if err := db.First(&studentRole, "name = ?", "student").Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Student role not found"})
	}

	// create user
	user := models.User{
		Email:  st.Email,
		Name:   st.Name,
		RoleID: studentRole.ID,
	}
	if err := user.SetPassword(st.Password); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Password hashing failed"})
	}
	if err := db.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "User creation failed", "details": err.Error()})
	}

	// create student
	student := models.Student{
		UserID:     user.ID,
		BatchID:    st.BatchID,
		RollNumber: st.RollNumber,
		Branch:     st.Branch,
	}

	if err := h.service.CreateStudent(&student); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Student creation failed"})
	}

	return c.Status(201).JSON(student)
}

// ==============================================
//
//	CREATE **MULTIPLE** STUDENTS  (/bulk)
//
// ==============================================
func (h *StudentHandler) CreateManyStudents(c *fiber.Ctx) error {
	var input CreateStudentBatchInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON body"})
	}

	if len(input.Students) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No students provided"})
	}

	db := c.Locals("db").(*gorm.DB)
	claims := c.Locals("claims").(*auth.Claims)

	if claims.Role != "institution" {
		return c.Status(403).JSON(fiber.Map{"error": "Only institutions can create students"})
	}

	var studentRole models.Role
	if err := db.First(&studentRole, "name = ?", "student").Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Student role not found"})
	}

	var createdStudents []models.Student

	for _, st := range input.Students {
		if st.BatchID == uuid.Nil {
			return c.Status(400).JSON(fiber.Map{"error": "BatchID is required"})
		}

		// validate batch
		var batch models.Batch
		if err := db.First(&batch, "id = ?", st.BatchID).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid Batch ID"})
		}

		// create user
		user := models.User{
			Email:  st.Email,
			Name:   st.Name,
			RoleID: studentRole.ID,
		}
		if err := user.SetPassword(st.Password); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Password hashing failed"})
		}
		if err := db.Create(&user).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "User creation failed", "details": err.Error()})
		}

		// create student
		student := models.Student{
			UserID:     user.ID,
			BatchID:    st.BatchID,
			RollNumber: st.RollNumber,
			Branch:     st.Branch,
		}
		if err := h.service.CreateStudent(&student); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Student creation failed"})
		}

		createdStudents = append(createdStudents, student)
	}

	return c.Status(201).JSON(fiber.Map{
		"message":  "Students created successfully",
		"total":    len(createdStudents),
		"students": createdStudents,
	})
}

// ==============================================
//
//	GET STUDENTS BY BATCH ID
//
// ==============================================
func (h *StudentHandler) GetByBatch(c *fiber.Ctx) error {
	batchID := c.Params("batchID")

	if batchID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Batch ID is required"})
	}

	students, err := h.service.GetByBatch(batchID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch students", "details": err.Error()})
	}

	return c.JSON(students)
}
