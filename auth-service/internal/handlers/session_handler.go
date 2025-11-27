package handlers

import (
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"auth-service/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SessionHandler struct {
	service *services.SessionService
}

func NewSessionHandler(db *gorm.DB) *SessionHandler {
	repo := repositories.NewSessionRepository(db)
	service := services.NewSessionService(repo)
	return &SessionHandler{service}
}

// 📌 1) Create Session (used in LOGIN)
func (h *SessionHandler) CreateSession(c *fiber.Ctx) error {
	var s models.Session
	if err := c.BodyParser(&s); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	if err := h.service.CreateSession(&s); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create session"})
	}

	return c.Status(201).JSON(s)
}

// 📌 2) Force Logout previous device (STUDENT ONLY)
func (h *SessionHandler) InvalidateOldSessions(c *fiber.Ctx) error {
	userID := c.Params("userID")

	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing userID"})
	}

	h.service.InvalidateOldSessions(userID)
	return c.JSON(fiber.Map{"success": true, "message": "Previous sessions invalidated"})
}

// 📌 3) Get All Active Sessions of User
func (h *SessionHandler) GetActiveSessions(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing userID"})
	}

	sessions, err := h.service.GetActiveSessions(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error fetching sessions"})
	}
	return c.JSON(sessions)
}

// 📌 4) Logout Specific Session  →  /session/logout/:token
func (h *SessionHandler) LogoutSession(c *fiber.Ctx) error {
	token := c.Params("token")
	userID := c.Params("userID")

	if token == "" || userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid params"})
	}

	err := h.service.DeactivateSession(userID, token)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Logout failed"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}
