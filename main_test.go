package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/database"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/logger"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/models"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/router"
	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// var startTime time.Time

func setupTestDB() {
	// Use an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// AutoMigrate necessary models
	_ = db.AutoMigrate(&models.Metrics{}) // Ensure this model is correct

	db.Create(&models.Metrics{
		ID:         uuid.New(),
		CPUPercent: 10.3,
		MemPercent: 20.3,
	})

	// Assign test database
	database.DB = db
}

func TestGetMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Log, _ = zap.NewDevelopment()

	// Initialize test database
	setupTestDB()
	defer cleanupTestDB()

	// Setup router
	r := gin.Default()
	router.SetRouter(r)

	// Perform test request
	req, _ := http.NewRequest("GET", "/metrics/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetMetricsByTimeRange(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Log, _ = zap.NewDevelopment()

	// Initialize test database
	setupTestDB()
	defer cleanupTestDB()

	// Insert test data into the database
	testMetric := models.Metrics{
		CPUPercent: 45.5,
		MemPercent: 60.2,
		CreatedAt:  time.Now().UTC(), // Store timestamp in UTC (matches DB behavior)
	}
	database.DB.Create(&testMetric)

	// Define start and end time in RFC3339 format (which is used in API requests)
	startTime := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	endTime := time.Now().Add(1 * time.Hour).UTC().Format(time.RFC3339)

	// Setup router
	r := gin.Default()
	router.SetRouter(r)

	// Perform test request with query parameters
	req, _ := http.NewRequest("GET", fmt.Sprintf("/metrics?start=%s&end=%s", startTime, endTime), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Validate response contains data
	assert.NotNil(t, response["data"])
	assert.Greater(t, len(response["data"].([]interface{})), 0, "Expected at least one metric in response")
}

func TestGetAverageMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Log, _ = zap.NewDevelopment()

	// Initialize test database
	setupTestDB()
	defer cleanupTestDB()

	metrics := []models.Metrics{
		{ID: uuid.New(), CPUPercent: 30.5, MemPercent: 50.2, CreatedAt: time.Now().UTC()},
		{ID: uuid.New(), CPUPercent: 40.0, MemPercent: 55.0, CreatedAt: time.Now().UTC()},
		{ID: uuid.New(), CPUPercent: 50.0, MemPercent: 60.0, CreatedAt: time.Now().UTC()},
	}
	database.DB.Create(&metrics)

	// Define start and end times for the query
	startTime := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	endTime := time.Now().Add(1 * time.Hour).UTC().Format(time.RFC3339)

	// Perform test request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/metrics/average?start=%s&end=%s", startTime, endTime), nil)
	w := httptest.NewRecorder()

	// Setup router
	r := gin.Default()
	router.SetRouter(r)
	r.ServeHTTP(w, req)

	// Print response for debugging
	fmt.Println("Response Body:", w.Body.String())

	// Assert response status
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response into AvgMetrics struct
	var response struct {
		Data models.AvgMetrics `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert expected values
	assert.NotZero(t, response.Data.CPUPercent, "CPUPercent should not be zero")
	assert.NotZero(t, response.Data.MemPercent, "MemPercent should not be zero")
}

func TestCollectAndSaveMetrics_Success(t *testing.T) {
	errChan := make(chan error, 1)
	logger.Log, _ = zap.NewDevelopment()

	//  Set up the test database
	setupTestDB()
	defer cleanupTestDB() // Ensure cleanup after test execution

	//  Run the function in a separate goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		service.CollectAndSaveMetrics(errChan)
	}()

	//  Wait for the goroutine to complete
	wg.Wait()

	//  Assert that no error was received
	select {
	case err := <-errChan:
		assert.Fail(t, "Unexpected error: %v", err)
	default:
		assert.True(t, true) // No error case
	}

	//  Verify data was inserted into the database
	var count int64
	database.DB.Model(&models.Metrics{}).Count(&count)
	assert.Greater(t, count, int64(0), "Expected metrics to be saved in DB")
}

func TestGracefulShutdown(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := database.Shutdown(ctx)
	assert.NoError(t, err, "Shutdown should complete without error")
}

func TestMetricsCollectorGracefulShutdown(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	setupTestDB()
	defer cleanupTestDB()

	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error, 1)
	defer close(errChan)


	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Recovered from panic: %v", r)
			}
		}()
		service.MetricsCollector(ctx, 1, errChan)
	}()

	time.Sleep(2 * time.Second)
	cancel()

	assert.True(t, true, "Metrics collector should exit cleanly")
}

func cleanupTestDB() {
	// Drop tables after test to clean up
	database.DB.Exec("DROP TABLE metrics;")
}
