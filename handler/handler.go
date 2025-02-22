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
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /metrics/ [get]
func GetAllMetrics(c *gin.Context) {
	logger.Log.Debug("GetAllMetrics handler")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// Calculate offset
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	response, totalRecords, err := service.GetAllMetrics(context.Background(), pageSize, offset)
	if err != nil {
		logger.Log.Error("GetAllMetric error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err,
			"time":    time.Now(),
		})
		return
	}

	if len(response) == 0 {
		logger.Log.Error("No metrics found ")
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No metrics found ",
			"time":    time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":         response,
		"totalRecords": totalRecords,
		"time":         time.Now(),
	})
}

// GetMetricsByTimeRange retrieves metrics within a specified time range.
//
// @Summary      Get Metrics By Time Range
// @Description  Fetches all metrics collected between the provided start and end timestamps.
// @Tags         Metrics
// @Accept       json
// @Produce      json
// @Param        start  query     string  true  "Start timestamp (RFC3339 format, e.g., 2025-02-22T00:00:00Z)"
// @Param        end    query     string  true  "End timestamp (RFC3339 format, e.g., 2025-02-22T23:59:59Z)"
// @Success      200    {object}  map[string]interface{}  "List of metrics"
// @Failure      400    {object}  map[string]interface{}  "Invalid input format"
// @Failure      404    {object}  map[string]interface{}  "No metrics found"
// @Failure      500    {object}  map[string]interface{}  "Internal server error"
// @Router       /metrics [get]
func GetMetricsByTimeRange(c *gin.Context) {

	logger.Log.Debug("GetMetricsByTimeRange handler")

	start := c.Query("start")
	end := c.Query("end")

	parsedStartTime, err := utils.ParseTime(start)
	if err != nil {
		logger.Log.Error("start time parsing error ", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "start time parsing error",
			"Error":   err,
			"time":    time.Now(),
		})
		return
	}

	parsedEndTime, err := utils.ParseTime(end)
	if err != nil {
		logger.Log.Error("End time parsing error ", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "end time parsing error",
			"Error":   err,
			"time":    time.Now(),
		})
		return
	}

	response, err := service.GetMetricsByTimeRange(context.Background(), parsedStartTime, parsedEndTime)
	if err != nil {
		logger.Log.Error("GetMetricsByTimeRange error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err,
			"time":    time.Now(),
		})
		return
	}

	if len(response) == 0 {
		logger.Log.Error("No metrics found in the given time range")
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No metrics found in the given time range",
			"time":    time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
		"time": time.Now(),
	})
}

// GetAverageMetrics godoc
// @Summary Get average CPU and memory usage in a time range
// @Description Retrieve average CPU and memory usage between start and end timestamps
// @Tags Metrics
// @Accept  json
// @Produce  json
// @Param start query string true "Start timestamp (RFC3339 format)"
// @Param end query string true "End timestamp (RFC3339 format)"
// @Success 200 {object} models.Metrics
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /metrics/average [get]
func GetAverageMetrics(c *gin.Context) {

	logger.Log.Debug("GetAlGetAverageMetrics handler")

	start := c.Query("start")
	end := c.Query("end")

	parsedStartTime, err := utils.ParseTime(start)
	if err != nil {
		logger.Log.Error("start time parsing error ", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "start time parsing error",
			"Error":   err,
			"time":    time.Now(),
		})
		return
	}

	parsedEndTime, err := utils.ParseTime(end)
	if err != nil {
		logger.Log.Error("End time parsing error ", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "end time parsing error",
			"Error":   err,
			"time":    time.Now(),
		})
		return
	}

	response, err := service.GetAverageMetrics(context.Background(), parsedStartTime, parsedEndTime)
	if err != nil {
		logger.Log.Error("GetAverageMetrics error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err,
			"time":  time.Now(),
		})
		return
	}

	if response.CPUPercent == 0 || response.MemPercent == 0 {
		logger.Log.Error("No metrics found in the given time range")
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No metrics found in the given time range",
			"time":    time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
		"time": time.Now(),
	})
}
