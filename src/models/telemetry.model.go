package models

import "time"

// Metric represents a telemetry metric.
type Metric struct {
	Name  string            `json:"name"`
	Value float64           `json:"value"`
	Tags  map[string]string `json:"tags"`
	Ts    time.Time         `json:"ts"`
}

// Event represents a telemetry event.
type Event struct {
	Name     string            `json:"name"`
	Metadata map[string]any    `json:"metadata"`
	Tags     map[string]string `json:"tags"`
	Ts       time.Time         `json:"ts"`
}
