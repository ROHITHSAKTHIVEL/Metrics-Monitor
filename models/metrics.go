package models

import (
	"time"

	"github.com/google/uuid"
)

type Metrics struct {
	ID         uuid.UUID `gorm:"primaryKey" json:"id"`
	CPUPercent float64   `gorm:"not null" json:"cpu_percent"`
	MemPercent float64   `gorm:"not null" json:"mem_percent"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type AvgMetrics struct {
	CPUPercent float64 `gorm:"not null" json:"cpu_percent"`
	MemPercent float64 `gorm:"not null" json:"mem_percent"`
}
