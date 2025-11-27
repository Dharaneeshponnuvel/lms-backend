package handlers

import (
	"auth-service/internal/auth"
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"auth-service/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

/* ------------ Create Institution + User (Login) ------------- */

type InstitutionRequest struct {
	Email         string         `json:"email"`
	Password      string         `json:"password"`
	Name          string         `json:"name"`
	AdminPosition string         `json:"admin_position"`
	Settings      datatypes.JSON `json:"settings"`
}

type InstitutionHandler struct {
	service *services.InstitutionService
}

func NewInstitutionHandler(db *gorm.DB) *InstitutionHandler {
	repo := repositories.NewInstitutionRepository(db)
	service := services.NewInstitutionService(repo)
	return &InstitutionHandler{service}
}

// ▶ CREATE INSTITUTION
func (h *InstitutionHandler) CreateInstitution(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)

	var body InstitutionRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if !auth.IsValidEmail(body.Email) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email format"})
	}
	if !auth.IsValidPassword(body.Password) {
		return c.Status(400).JSON(fiber.Map{"error": "Weak password"})
	}

	var role models.Role
	if err := db.Where("name = ?", "INSTITUTION_ADMIN").First(&role).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Role not found"})
	}

	tx := db.Begin()

	user := models.User{
		Email:  body.Email,
		Name:   body.Name,
		RoleID: role.ID,
	}
	user.SetPassword(body.Password)

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return c.Status(409).JSON(fiber.Map{"error": "Email already exists"})
	}

	inst := models.Institution{
		UserID:        user.ID,
		Name:          body.Name,
		AdminPosition: body.AdminPosition,
		Settings:      body.Settings,
	}
	if err := tx.Create(&inst).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create institution"})
	}

	tx.Commit()

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"user":    user,
		"institution": fiber.Map{
			"id":     inst.ID,
			"name":   inst.Name,
			"userId": inst.UserID,
		},
	})
}

// ▶ GET INSTITUTION BY ID  (🚀 FIXED ERROR)
func (h *InstitutionHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	inst, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Institution not found"})
	}

	return c.JSON(inst)
}
