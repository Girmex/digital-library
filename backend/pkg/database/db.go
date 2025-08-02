package database

import (
	"digital-library/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=1234 dbname=digital_library port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	
	DB = db
	
	// Auto migrate models
	err = DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Book{}, &models.Category{}, &models.Loan{})
	if err != nil {
		log.Fatal("Failed to migrate database")
	}
}