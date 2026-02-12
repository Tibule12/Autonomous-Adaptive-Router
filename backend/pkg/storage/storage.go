package storage

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type MetricPoint struct {
	Time       string  `json:"time"`
	Latency    int64   `json:"latency"`
	PacketLoss float64 `json:"packet_loss"`
}

var (
	historyFile = "metrics_history.json"
	mu          sync.Mutex
)

// SaveMetrics appends a new metric point to the history file
func SaveMetric(latency int64, packetLoss float64) error {
	mu.Lock()
	defer mu.Unlock()

	// Load existing
	var points []MetricPoint
	data, err := os.ReadFile(historyFile)
	if err == nil {
		json.Unmarshal(data, &points)
	}

	// Add new point
	points = append(points, MetricPoint{
		Time:       time.Now().Format("15:04:05"),
		Latency:    latency,
		PacketLoss: packetLoss,
	})

	// Keep last 50 points
	if len(points) > 50 {
		points = points[len(points)-50:]
	}

	// Save back
	newData, _ := json.MarshalIndent(points, "", "  ")
	return os.WriteFile(historyFile, newData, 0644)
}

// GetHistory returns the stored metrics
func GetHistory() []MetricPoint {
	mu.Lock()
	defer mu.Unlock()

	var points []MetricPoint
	data, err := os.ReadFile(historyFile)
	if err != nil {
		return points // Return empty if no file
	}
	json.Unmarshal(data, &points)
	return points
}
