package service

import (
	"time"

	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/database"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/logger"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/models"
	"github.com/google/uuid"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"
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

	// Get CPU Usage
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		logger.Log.Error("Failed to get CPU usage", zap.Error(err))
		errChan <- err
		return
	}
	metrics.CPUPercent = cpuPercent[0]

	logger.Log.Info("CPU Percent", zap.Float64("value", cpuPercent[0]))

	// Get Memory Usage
	memStats, err := mem.VirtualMemory()
	if err != nil {
		logger.Log.Error("Failed to get Memory usage", zap.Error(err))
		errChan <- err
		return
	}
	metrics.MemPercent = memStats.UsedPercent

	logger.Log.Info("Memory Percent", zap.Float64("value", memStats.UsedPercent))

	// Store in database
	if err := database.DB.Create(&metrics).Error; err != nil {
		logger.Log.Error("Failed to insert metrics into database", zap.Error(err))
		errChan <- err
	}
}
