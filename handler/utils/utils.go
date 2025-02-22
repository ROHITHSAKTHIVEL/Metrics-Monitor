package utils

import (
	"time"

	"github.com/ROHITHSAKTHIVEL/Metrics-Monitor/logger"
	"go.uber.org/zap"
)

func ParseTime(req string) (time.Time, error) {

	parsedTime, err := time.Parse(time.RFC3339, req)
	if err != nil {
		logger.Log.Error("Invalid  timestamp:", zap.Error(err))
		return time.Time{}, err
	}
	return parsedTime, nil

}
