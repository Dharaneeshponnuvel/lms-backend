package database

import (
	"fmt"
	"log"

	"auth-service/internal/config"
	"auth-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("DB connection failed: %v", err)
	}

	log.Println("📦 Database connected!")

	db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Institution{},
		&models.BatchYear{},
		&models.Batch{},
		&models.Student{},
		&models.ContentManager{},
		&models.ContentManagerBatch{},
		&models.Session{},
		&models.RefreshToken{},
	)

	SeedRoles(db) // ⬅ MUST CALL

	return db, nil
}
