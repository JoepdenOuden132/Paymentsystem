package database

import (
	"log"
	"os"

	"main.go/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("fonteyn:#Funckypower1.@tpc(db-gen-01.mysql.database.azure.com:3306)/restapi")
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	DB.AutoMigrate(&models.Payment{})
}
