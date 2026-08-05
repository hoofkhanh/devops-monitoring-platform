package model

import (
	"fmt"
	"time"
)

type Metric struct {
	ID        int64     `json:"id"`
	CPU       float64   `json:"cpu"`
	Memory    float64   `json:"memory"`
	Disk      float64   `json:"disk"`
	Timestamp time.Time `json:"timestamp"`
}

type MetricCreateRequest struct {
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
	Disk   float64 `json:"disk"`
}

func (m MetricCreateRequest) Validate() error {
	if m.CPU < 0 || m.Memory < 0 || m.Disk < 0 {
		return fmt.Errorf("metric values must be non-negative")
	}
	return nil
}
