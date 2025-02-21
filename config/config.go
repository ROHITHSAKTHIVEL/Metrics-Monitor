package config

import (
	"os"
	"strconv"

	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/logger"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/models"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func LoadConfig() *models.Config {
	err := godotenv.Load()
	if err != nil {
		logger.Log.Fatal("Error In Loading Config", zap.Error(err))
	}

	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	metricsInterval, _ := strconv.Atoi(os.Getenv("METRICS_INTERVAL_SECONDS"))

	return &models.Config{
		DBHost:          os.Getenv("DB_HOST"),
		DBUser:          os.Getenv("DB_USER"),
		DBPass:          os.Getenv("DB_PASS"),
		DBName:          os.Getenv("DB_NAME"),
		Port:            os.Getenv("PORT"),
		DBPort:          dbPort,
		MetricsInterval: metricsInterval,
	}

}
