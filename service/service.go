package service

import (
	"context"
	"time"

	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/database"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/logger"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/models"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/utils"
	"github.com/google/uuid"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func MetricsCollector(interval int, stopChan chan struct{}, errChan chan error) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logger.Log.Info("Collecting system metrics...")
			go CollectAndSaveMetrics(errChan)

		case <-stopChan:
			logger.Log.Info("Stopping Metrics Collector...")
			return

		case err := <-errChan:
			logger.Log.Error("Error in CollectAndSaveMetrics", zap.Error(err))
			continue
		}
	}
}

// CollectAndSaveMetrics collects CPU and memory metrics and saves them to DB
func CollectAndSaveMetrics(errChan chan error) {
	var metrics models.Metrics
	metrics.ID = uuid.New()

	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		logger.Log.Error("Failed to get CPU usage", zap.Error(err))
		errChan <- err
		return
	}
	metrics.CPUPercent = utils.RoundToTwoDecimal(cpuPercent[0])

	logger.Log.Info("CPU Percent", zap.Float64("value", cpuPercent[0]))

	memStats, err := mem.VirtualMemory()
	if err != nil {
		logger.Log.Error("Failed to get Memory usage", zap.Error(err))
		errChan <- err
		return
	}
	metrics.MemPercent = utils.RoundToTwoDecimal(memStats.UsedPercent)

	logger.Log.Info("Memory Percent", zap.Float64("value", memStats.UsedPercent))

	if err := database.DB.Create(&metrics).Error; err != nil {
		logger.Log.Error("Failed to insert metrics into database", zap.Error(err))
		errChan <- err
	}
}

func GetAllMetrics(ctx context.Context, pageSize, offset int) ([]models.Metrics, int64, error) {
	var res []models.Metrics
	var totalRecords int64

	if database.DB == nil {
		logger.Log.Error("Database is not initialized")
		return nil, 0, gorm.ErrInvalidDB
	}

	if err := database.DB.Model(&models.Metrics{}).Count(&totalRecords).Error; err != nil {
		logger.Log.Error("Error counting metrics:", zap.Error(err))
		return nil, 0, err
	}

	if err := database.DB.WithContext(ctx).
		Limit(pageSize).
		Offset(offset).
		Order("created_at DESC").
		Find(&res).Error; err != nil {
		logger.Log.Error("Error fetching metrics:", zap.Error(err))
		return nil, 0, err
	}

	return res, totalRecords, nil
}

func GetMetricsByTimeRange(ctx context.Context, start, end time.Time) ([]models.Metrics, error) {

	var res []models.Metrics

	err := database.DB.WithContext(ctx).
		Where("created_at BETWEEN ? AND ?", start, end).
		Find(&res).Error
	if err != nil {
		logger.Log.Error("Error fetching metrics by time range:", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func GetAverageMetrics(ctx context.Context, start, end time.Time) (models.AvgMetrics, error) {
	var res models.AvgMetrics

	if database.DB == nil {
		logger.Log.Error("Database is not initialized")
		return models.AvgMetrics{}, gorm.ErrInvalidDB
	}

	if err := database.DB.WithContext(ctx).
		Raw(`SELECT 
                AVG(cpu_percent) AS cpu_percent, 
                AVG(mem_percent) AS mem_percent 
             FROM metrics 
             WHERE created_at BETWEEN ? AND ?`, start, end).
		Scan(&res).Error; err != nil {
		logger.Log.Error("Error fetching average metrics:", zap.Error(err))
		return models.AvgMetrics{}, err
	}

	return res, nil
}
