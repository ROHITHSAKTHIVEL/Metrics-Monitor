package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "github.com/lib/pq"
)

var DB *gorm.DB

func InitDB(cfg *models.Config) {

	dbURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)

	CeateDbNotExist(dbURL, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("error in connecting DB")
		return
	}
	DB = db
	DB.AutoMigrate(
		&models.Metrics{},
	)
}

func CeateDbNotExist(dburl string, dbName string) {
	sqlDB, err := sql.Open("postgres", dburl)
	if err != nil {
		log.Fatalf("error connecting to default database: %v", err)
		return
	}
	defer sqlDB.Close()

	// Check if the database exists
	var exists bool
	err = sqlDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		log.Fatalf("error checking database existence: %v", err)
		return
	}

	// Create the database if it does not exist
	if !exists {
		_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			log.Fatalf("error creating database: %v", err)
			return
		}
		log.Printf("Database %s created successfully!", dbName)
	}
}
