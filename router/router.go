package router

import (
	_ "github.com/ROHITHSAKTHIVEL/Metrics-Monitor/docs"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/handler"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetRouter(apiRouter *gin.Engine) {

	metrics := apiRouter.Group("/metrics")
	{
		metrics.GET("/", handler.GetAllMetrics)
		metrics.GET("", handler.GetMetricsByTimeRange)
		metrics.GET("/average", handler.GetAverageMetrics)
	}
	apiRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
