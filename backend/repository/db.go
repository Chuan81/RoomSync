package repository

import (
	"log"
	"roomsync/config"
	"roomsync/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := config.AppConfig.Database.DSN
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto Migrate
	err = DB.AutoMigrate(&models.User{}, &models.Room{}, &models.Booking{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
