package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/config"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/database"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/logger"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/router"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	cfg := config.LoadConfig()
	database.InitDB(cfg)
	r := gin.Default()
	router.SetRouter(r)

	stopChan := make(chan struct{})
	errChan := make(chan error)
	go service.MetricsCollector(cfg.MetricsInterval, stopChan, errChan)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		for err := range errChan {
			logger.Log.Error("Metrics collection error", zap.Error(err))
		}
	}()

	go func() {
		if err := r.Run(cfg.Port); err != nil {
			logger.Log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for termination signal
	sig := <-sigChan
	logger.Log.Info("Received termination signal", zap.String("signal", sig.String()))

	// Stop Metrics Collector
	close(stopChan)
	logger.Log.Info("Shutting down gracefully...")
	time.Sleep(2 * time.Second) // Give time for goroutines to clean up
}
