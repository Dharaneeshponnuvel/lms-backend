package handlers

import (
	"regexp"
	"strings"
	"time"

	"auth-service/internal/auth"
	"auth-service/internal/config"
	"auth-service/internal/models"
	"auth-service/internal/redis"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuthHandler struct{ cfg *config.Config }

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

/* -------------------- 🔐 PASSWORD RULES -------------------- */

// ---------------- PASSWORD CHECK ----------------
func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	lower := regexp.MustCompile(`[a-z]`)
	upper := regexp.MustCompile(`[A-Z]`)
	number := regexp.MustCompile(`[0-9]`)
	special := regexp.MustCompile(`[@$!%*?&]`)

	return lower.MatchString(password) &&
		upper.MatchString(password) &&
		number.MatchString(password) &&
		special.MatchString(password)
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

/* -------------------- REGISTER -------------------- */

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)

	var body registerRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	if !isValidEmail(body.Email) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email format"})
	}
	if !isValidPassword(body.Password) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Password must be 8+ characters & include UPPER, lower, number & special character",
		})
	}

	var role models.Role
	if err := db.Where("name = ?", body.Role).First(&role).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid role name"})
	}

	user := models.User{Email: body.Email, Name: body.Name, RoleID: role.ID}
	if err := user.SetPassword(body.Password); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "Email already exists"})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "user": user})
}

/* -------------------- LOGIN -------------------- */

type loginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	DeviceInfo string `json:"device_info"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)

	var body loginRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	var user models.User
	if err := db.Preload("Role").Where("email = ?", body.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	if !user.CheckPassword(body.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if user.Role.Name == "STUDENT" {
		db.Model(&models.Session{}).Where("user_id = ? AND is_active = true", user.ID).
			Updates(map[string]interface{}{"is_active": false})
		redis.RDB.Del(redis.Ctx, "session:"+user.ID.String())
	}

	accessToken, err := auth.GenerateAccessToken(h.cfg, user.ID.String(), user.Email, user.Role.Name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate access token"})
	}

	// Create Session
	session := models.Session{
		UserID:     user.ID,
		Token:      accessToken,
		DeviceInfo: datatypes.JSON([]byte(body.DeviceInfo)),
		ExpiresAt:  time.Now().Add(time.Hour * 1),
		IsActive:   true,
	}
	db.Create(&session)

	// Create Refresh Token
	refreshToken := models.RefreshToken{
		UserID:    user.ID,
		SessionID: session.ID,
		Token:     auth.GenerateSecureToken(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}
	db.Create(&refreshToken)

	redis.RDB.Set(redis.Ctx, "session:"+user.ID.String(), accessToken, time.Hour*1)

	return c.JSON(fiber.Map{
		"success":       true,
		"access_token":  accessToken,
		"refresh_token": refreshToken.Token,
		"user": fiber.Map{
			"id":     user.ID,
			"email":  user.Email,
			"name":   user.Name,
			"role":   user.Role.Name,
			"roleId": user.RoleID,
		},
	})
}

/* -------------------- REFRESH TOKEN -------------------- */

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)

	refresh := c.FormValue("refresh_token")
	if refresh == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing refresh_token"})
	}

	var token models.RefreshToken
	if err := db.Where("token = ? AND is_used = false", refresh).First(&token).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired refresh token"})
	}

	db.Model(&models.RefreshToken{}).Where("id = ?", token.ID).Update("is_used", true)

	var user models.User
	db.First(&user, "id = ?", token.UserID)

	accessToken, err := auth.GenerateAccessToken(h.cfg, user.ID.String(), user.Email, user.Role.Name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Token generation failed"})
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"access_token": accessToken,
	})
}

/* -------------------- LOGOUT -------------------- */

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*auth.Claims)
	db := c.Locals("db").(*gorm.DB)

	redis.RDB.Del(redis.Ctx, "session:"+claims.UserID)

	db.Model(&models.Session{}).
		Where("user_id = ? AND token = ?", claims.UserID, extractToken(c.Get("Authorization"))).
		Update("is_active", false)

	return c.JSON(fiber.Map{"success": true, "message": "Logged out"})
}

func (h *AuthHandler) Verify(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*auth.Claims)
	return c.JSON(fiber.Map{"success": true, "user": claims})
}

/* -------------------- UTILITY -------------------- */
func extractToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.Split(header, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
