package router

import (
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/handler"
	"github.com/gin-gonic/gin"
)

func SetRouter(apiRouter *gin.Engine) {

	metrics := apiRouter.Group("/metrics")
	{
		metrics.GET("/", handler.GetAllMetrics)
		metrics.GET("", handler.GetMetricsByTimeRange)
		metrics.GET("/average", handler.GetAverageMetrics)
	}

}
