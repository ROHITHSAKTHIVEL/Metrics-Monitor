package main

import (
	"context"
	"net/http"
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

	server := &http.Server{
		Addr:    cfg.Port,
		Handler: r,
	}

	ctx, cancel := context.WithCancel(context.Background())

	errChan := make(chan error, 10)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go service.MetricsCollector(ctx, cfg.MetricsInterval, errChan)

	go func() {
		for err := range errChan {
			logger.Log.Error("Metrics collection error", zap.Error(err))
		}
	}()

	logger.Log.Info("Starting API server", zap.String("port", cfg.Port))
	// Used standard http lib for graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Wait for termination signal
	sig := <-sigChan
	logger.Log.Info("Received termination signal", zap.String("signal", sig.String()))

	// Stop Metrics Collector gracefully
	cancel()
	time.Sleep(1 * time.Second)
	close(errChan)

	// Close Database
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer shutdownCancel()

	//shutdown http server gracefully
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("Server shutdown failed", zap.Error(err))
	} else {
		logger.Log.Info("Server shutdown successfully")
	}

	// Close database connection
	if err := database.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("Database shutdown failed", zap.Error(err))
	} else {
		logger.Log.Info("Database shutdown successfully")
	}

	logger.Log.Info("Shutting down gracefully...")
	time.Sleep(1 * time.Second)
	os.Exit(0)
}
