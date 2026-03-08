package database

import (
	"fmt"
	"os"

	"github.com/Dukvaha27/film-library/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("ошибка при инициализации DB: %w", err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Movie{},
		&models.Watchlist{},
		&models.Review{},
		&models.Genre{},
	)
	if err != nil {
		return fmt.Errorf("ошибка при миграции таблиц: %w", err)
	}

	return nil
}
