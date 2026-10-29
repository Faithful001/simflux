package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("simflux.db"), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		DB = nil
		return
	}
}

func GetDB() *gorm.DB {
	return DB
}

func AutoMigrate(models ...interface{}) error {
	return DB.AutoMigrate(models...)
}