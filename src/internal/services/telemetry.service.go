package services

import (
	"context"
)

// TelemetryService defines the interface for storing and retrieving telemetry data.
type TelemetryService interface {
	// RecordEvent logs a specific event with associated metadata.
	RecordEvent(ctx context.Context, eventName string, metadata map[string]any, tags map[string]string) error

	// RecordMetric logs a numerical metric with a value and optional tags.
	RecordMetric(ctx context.Context, metricName string, value float64, tags map[string]string) error

	// Upload uploads telemetry the enqueued data. This method should be called
	// in a separate goroutine to avoid blocking. Should panic if impl is not async.
	Upload() error

	// Creates a counter metric with the given name and tags.
	GetCounter(ctx context.Context, metricName string, tags map[string]string) (Counter, error)
}

// Counter defines the interface for a counter metric.
type Counter interface {
	// Increment increments the counter by the given count.
	Increment(count uint64)

	// Upload uploads telemetry the enqueued data. This method should be called
	// in a separate goroutine to avoid blocking. Should panic if impl is not async.
	Upload() error
}
