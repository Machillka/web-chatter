package db

import (
	"github.com/machillka/web-chatter/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	dsn := config.MySQLDSN()
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移 User 模型
	err = gormDB.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	DB = gormDB
	return nil
}
