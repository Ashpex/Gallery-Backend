package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gitlab.zalopay.vn/top/intern/vybnt/gallery-backend/gallery/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

// Config is a struct for config configuration
// SetupDB is a function to setup config connection
func SetupDB() *gorm.DB {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Failed to load env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: false})

	if err != nil {
		panic("Failed to create a connection to config")
	}
	db.AutoMigrate(&entity.User{}, &entity.Post{}, &entity.Topic{}, &entity.Comment{}, &entity.Follower{}, &entity.Like{}, &entity.Subscribe{})
	return db
}

// CloseDB is a function to close config connection
func CloseDB(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		fmt.Sprintf("Failed to close connection to config")
	}
	dbSQL.Close()
}
