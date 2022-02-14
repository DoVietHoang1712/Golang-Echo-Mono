package db

import (
	"fmt"
	"golang-sql/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New() *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:17122000@tcp(127.0.0.1:3306)/db?charset=utf8&parseTime=True&loc=Local",
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(3)
	db.Logger.LogMode(logger.LogLevel(4))
	return db
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.Article{},
		&model.User{},
		&model.Follow{},
		&model.Tag{},
		&model.Comment{},
	)
}
