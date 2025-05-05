package config

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("seaguard.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info), // 添加更详细的日志
	})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// // 自动迁移表结构
	// err = DB.AutoMigrate(&models.Activity{}, &models.Volunteer{}, &models.Registration{})
	// if err != nil {
	// 	log.Fatal("Failed to migrate database:", err)
	// }
}
