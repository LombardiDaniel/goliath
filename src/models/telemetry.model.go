package models

import "time"

// Metric represents a telemetry metric.
type Metric struct {
	Name  string            `json:"name" bson:"name"`
	Tags  map[string]string `json:"tags" bson:"tags"`
	Value float64           `json:"value" bson:"value"`
	Ts    time.Time         `json:"ts" bson:"ts"`
}

// Event represents a telemetry event.
type Event struct {
	Name     string            `json:"name" bson:"name"`
	Tags     map[string]string `json:"tags" bson:"tags"`
	Metadata map[string]any    `json:"metadata" bson:"metadata"`
	Ts       time.Time         `json:"ts" bson:"ts"`
}
