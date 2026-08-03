package model

import (
	"fmt"
	"net"
	"time"
)

const (
	ServerStatusUnknown = "unknown"
	ServerStatusOnline  = "online"
	ServerStatusOffline = "offline"
)

type Server struct {
	ID        int64      `json:"id"`
	Hostname  string     `json:"hostname"`
	IP        string     `json:"ip"`
	OS        string     `json:"os"`
	Status    string     `json:"status"`
	LastSeen  *time.Time `json:"last_seen,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type ServerRegisterRequest struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	OS       string `json:"os"`
}

func (r ServerRegisterRequest) Validate() error {
	if r.Hostname == "" {
		return fmt.Errorf("hostname is required")
	}
	if r.IP == "" {
		return fmt.Errorf("ip is required")
	}
	if net.ParseIP(r.IP) == nil {
		return fmt.Errorf("ip must be a valid IP address")
	}
	if r.OS == "" {
		return fmt.Errorf("os is required")
	}
	return nil
}

type Metric struct {
	ID        int64     `json:"id"`
	ServerID  int64     `json:"server_id"`
	CPU       float64   `json:"cpu"`
	Memory    float64   `json:"memory"`
	Disk      float64   `json:"disk"`
	Timestamp time.Time `json:"timestamp"`
}

type MetricCreateRequest struct {
	ServerID int64   `json:"server_id"`
	CPU      float64 `json:"cpu"`
	Memory   float64 `json:"memory"`
	Disk     float64 `json:"disk"`
}

func (m MetricCreateRequest) Validate() error {
	if m.ServerID <= 0 {
		return fmt.Errorf("server_id must be greater than zero")
	}
	if m.CPU < 0 || m.Memory < 0 || m.Disk < 0 {
		return fmt.Errorf("metric values must be non-negative")
	}
	return nil
}
