package database

import (
	"log"

	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *models.Config) {
	dsn := "host=localhost user=postgres password=postgres dbname=employe port=5432 sslmode=disable "
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("error in connecting DB")
		return
	}
	DB = db
	DB.AutoMigrate(
		&models.Metrics{},
	)
}
