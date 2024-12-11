package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := "host=localhost user=roshan password=ASHIRWAD_1 dbname=airbnb port=5432 sslmode=disable"
	// DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to connect to database!")
	}
	fmt.Println("Database connected")
}
