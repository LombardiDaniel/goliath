package services

import (
	"context"

	"github.com/LombardiDaniel/gopherbase/models"
)

// TelemetryService defines the interface for storing and retrieving telemetry data.
type TelemetryService interface {
	// RecordEvent logs a specific event with associated metadata.
	RecordEvent(ctx context.Context, eventName string, metadata map[string]any) error

	// RecordMetric logs a numerical metric with a value and optional tags.
	RecordMetric(ctx context.Context, metricName string, value float64, tags map[string]string) error

	// RecordError logs an error with a message and optional metadata.
	RecordError(ctx context.Context, err error, metadata map[string]any) error

	// GetMetrics retrieves metrics based on a query (e.g., metric name, tags, time range).
	GetMetrics(ctx context.Context, query map[string]string) ([]models.Metric, error)

	// GetEvents retrieves events based on a query (e.g., event name, time range).
	GetEvents(ctx context.Context, query map[string]string) ([]models.Event, error)
}
