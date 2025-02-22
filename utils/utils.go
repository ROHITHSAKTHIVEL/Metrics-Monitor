package utils

import (
	"math"
	"time"

	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/logger"
	"go.uber.org/zap"
)

func ParseTime(req string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, req)
	if err != nil {
		logger.Log.Error("Invalid timestamp:", zap.Error(err))
		return time.Time{}, err
	}

	loc, _ := time.LoadLocation("Asia/Kolkata")
	istTime := time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(),
		parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(), loc)

	return istTime, nil
}

func RoundToTwoDecimal(num float64) float64 {
	return math.Round(num*100) / 100
}
