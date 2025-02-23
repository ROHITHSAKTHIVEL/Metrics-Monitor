package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/logger"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/models"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *models.Config) {

	dbURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)

	CeateDbNotExist(dbURL, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("error in connecting DB", zap.Error(err))
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
		logger.Log.Fatal("error connecting to default database: ", zap.Error(err))
		return
	}
	defer sqlDB.Close()

	// Check if the database exists
	var exists bool
	err = sqlDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		logger.Log.Fatal("error checking database existence:", zap.Error(err))
		return
	}

	// Create the database if it does not exist
	if !exists {
		_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			logger.Log.Fatal("error creating database: ", zap.Error(err))
			return
		}
		logger.Log.Info("Database created sucessfully ")
	}
}

func Shutdown(ctx context.Context) error {
	if DB == nil {
		logger.Log.Warn("Database is not initialized, skipping shutdown")
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		logger.Log.Error("Failed to get database connection", zap.Error(err))
		return err
	}

	// Use a channel to signal completion
	done := make(chan error, 1)

	go func() {
		logger.Log.Info("Closing database connection...")
		done <- sqlDB.Close()
	}()

	select {
	case <-ctx.Done():
		logger.Log.Warn("Shutdown timeout reached, force closing database connection")
		return ctx.Err()
	case err := <-done:
		if err != nil {
			logger.Log.Error("Error closing database connection", zap.Error(err))
			return err
		}
		logger.Log.Info("Database connection closed successfully")
		return nil
	}
}

