package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/logger"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/service"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetAllMetrics godoc
// @Summary Retrieve all collected metrics
// @Description Get paginated metrics data
// @Tags Metrics
// @Accept  json
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {array} models.Metrics
// @Failure 500 {object} map[string]string
// @Router /metrics/ [get]
func GetAllMetrics(c *gin.Context) {
	logger.Log.Debug("GetAllMetrics handler")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// Validate pagination inputs
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // Set max limit
	}
	offset := (page - 1) * pageSize

	response, totalRecords, err := service.GetAllMetrics(context.Background(), pageSize, offset)
	if err != nil {
		logger.Log.Error("GetAllMetrics error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err.Error(),
			"time":    time.Now().UTC(),
		})
		return
	}

	if len(response) == 0 {
		logger.Log.Warn("No metrics found")
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No metrics found",
			"time":    time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":         response,
		"totalRecords": totalRecords,
		"time":         time.Now().UTC(),
	})
}

// GetMetricsByTimeRange godoc
// @Summary Get Metrics By Time Range
// @Description Fetch metrics collected between start and end timestamps.
// @Tags Metrics
// @Accept json
// @Produce json
// @Param start query string true "Start timestamp (RFC3339 format, e.g., 2025-02-22T00:00:00Z)"
// @Param end query string true "End timestamp (RFC3339 format, e.g., 2025-02-22T23:59:59Z)"
// @Success 200 {array} models.Metrics
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /metrics [get]
func GetMetricsByTimeRange(c *gin.Context) {
	logger.Log.Debug("GetMetricsByTimeRange handler")

	start := c.Query("start")
	end := c.Query("end")

	parsedStartTime, err := utils.ParseTime(start)
	if err != nil {
		logger.Log.Error("Start time parsing error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid start time format",
			"Error":   err.Error(),
			"time":    time.Now().UTC(),
		})
		return
	}

	parsedEndTime, err := utils.ParseTime(end)
	if err != nil {
		logger.Log.Error("End time parsing error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid end time format",
			"Error":   err.Error(),
			"time":    time.Now().UTC(),
		})
		return
	}

	response, err := service.GetMetricsByTimeRange(context.Background(), parsedStartTime, parsedEndTime)
	if err != nil {
		logger.Log.Error("GetMetricsByTimeRange error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err.Error(),
			"time":    time.Now().UTC(),
		})
		return
	}

	if len(response) == 0 {
		logger.Log.Warn("No metrics found in the given time range")
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No metrics found in the given time range",
			"time":    time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
		"time": time.Now().UTC(),
	})
}

// GetAverageMetrics godoc
// @Summary Get average CPU and memory usage in a time range
// @Description Retrieve average CPU and memory usage between start and end timestamps
// @Tags Metrics
// @Accept json
// @Produce json
// @Param start query string true "Start timestamp (RFC3339 format)"
// @Param end query string true "End timestamp (RFC3339 format)"
// @Success 200 {object} models.AvgMetrics
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /metrics/average [get]
func GetAverageMetrics(c *gin.Context) {
	logger.Log.Debug("GetAverageMetrics handler")

	start := c.Query("start")
	end := c.Query("end")

	parsedStartTime, err := utils.ParseTime(start)
	if err != nil {
		logger.Log.Error("Start time parsing error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid start time format",
			"Error":   err.Error(),
			"time":    time.Now().UTC(),
		})
		return
	}

	parsedEndTime, err := utils.ParseTime(end)
	if err != nil {
		logger.Log.Error("End time parsing error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid end time format",
			"Error":   err.Error(),
			"time":    time.Now().UTC(),
		})
		return
	}

	response, err := service.GetAverageMetrics(context.Background(), parsedStartTime, parsedEndTime)
	if err != nil {
		logger.Log.Error("GetAverageMetrics error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
			"time":  time.Now().UTC(),
		})
		return
	}

	if response.CPUPercent == 0 && response.MemPercent == 0 {
		logger.Log.Warn("No metrics found in the given time range")
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No metrics found in the given time range",
			"time":    time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
		"time": time.Now().UTC(),
	})
}

// HealthCheck godoc
// @Summary Check service health
// @Description Returns service status
// @Tags Health
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
