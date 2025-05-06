package config

import (
	"log"
	"seaguard-admin-backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("seaguard.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

// 自动迁移表结构
err = DB.AutoMigrate(&models.User{}, &models.Activity{}, &models.Volunteer{}, &models.Registration{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
