package service

import (
	"ai-mockup/internal/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitPostgres() error {
	var err error
	DB, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		return err
	}
	// Auto-migrate the User model
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	return nil
}

func ClosePostgres() {
	if DB != nil {
		db, err := DB.DB()
		if err == nil {
			db.Close()
		}
	}
}
