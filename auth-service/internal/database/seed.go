package database

import (
	"auth-service/internal/models"
	"log"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	roles := []models.Role{
		{Name: "ADMIN"},
		{Name: "STUDENT"},
		{Name: "TEACHER"},
		{Name: "INSTITUTION"},
		{Name: "CONTENT_MANAGER"},
	}

	for _, role := range roles {
		db.FirstOrCreate(&role, models.Role{Name: role.Name})
		log.Println("✔ Role ensured in DB:", role.Name)
	}
}
