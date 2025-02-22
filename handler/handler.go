package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/handler/utils"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/logger"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetAllMetrics(c *gin.Context) {
	logger.Log.Debug("GetAllMetrics handler")

	response, err := service.GetAllMetrics(context.Background())
	if err != nil {
		logger.Log.Error("GetAllMetric error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err,
			"time":    time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": response,
		"time":    time.Now(),
	})
}

func GetMetricsByTimeRange(c *gin.Context) {

	logger.Log.Debug("GetMetricsByTimeRange handler")

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

	response, err := service.GetMetricsByTimeRange(context.Background(), parsedStartTime, parsedEndTime)
	if err != nil {
		logger.Log.Error("GetAllMetric error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err,
			"time":    time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": response,
		"time":    time.Now(),
	})
}

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
		logger.Log.Error("GetAllMetric error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err,
			"time":  time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": response,
		"time":    time.Now(),
	})
}
