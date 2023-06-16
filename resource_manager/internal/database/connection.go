package database

import (
	"log"
	"os"

	"resource_manager/internal/database/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

// Initializes a database connection
func Init() {
	dsn := os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_DATABASE") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DBConn = db

	db.AutoMigrate(
		&model.ScalingLog{},
	)
}
