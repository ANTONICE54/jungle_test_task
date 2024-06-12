package database

import (
	"app/internal/config"
	"app/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// Initialize connection to DB and runs migration
func InitDB(config *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Can`t establish connection with database:", err)

	}
	if exists := db.Migrator().HasTable(&models.User{}); !exists {
		db.Migrator().CreateTable(&models.User{})
	}

	if exists := db.Migrator().HasTable(&models.ImageURLs{}); !exists {
		db.Migrator().CreateTable(&models.ImageURLs{})
	}

	return db
}
