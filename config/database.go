package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/nabilulilalbab/welcomesite/models"
)

var DB *gorm.DB

func InitDatabase() {
	db, err := gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Task{})
	DB = db
}
