package services

import (
	"context"
	"time"

	"github.com/LombardiDaniel/gopherbase/common"
	"github.com/LombardiDaniel/gopherbase/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type TelemetryServiceMongoAsyncImpl struct {
	mongoClient     *mongo.Client
	metricsCol      *mongo.Collection
	eventsCol       *mongo.Collection
	batchInsertSize uint32
	metricCh        chan models.Metric
	eventsCh        chan models.Event
}

func NewTelemetryServiceMongoAsyncImpl(mongoClient *mongo.Client, metricsCol, eventsCol *mongo.Collection, batchInsertSize uint32) TelemetryService {
	return &TelemetryServiceMongoAsyncImpl{
		mongoClient:     mongoClient,
		metricsCol:      metricsCol,
		eventsCol:       eventsCol,
		metricCh:        make(chan models.Metric),
		eventsCh:        make(chan models.Event),
		batchInsertSize: batchInsertSize,
	}
}

func (s *TelemetryServiceMongoAsyncImpl) RecordEvent(ctx context.Context, eventName string, metadata map[string]any, tags map[string]string) error {
	e := models.Event{
		Name:     eventName,
		Metadata: metadata,
		Tags:     tags,
		Ts:       time.Now(),
	}
	s.eventsCh <- e
	return nil
}

func (s *TelemetryServiceMongoAsyncImpl) RecordMetric(ctx context.Context, metricName string, value float64, tags map[string]string) error {
	e := models.Metric{
		Name:  metricName,
		Value: value,
		Tags:  tags,
		Ts:    time.Now(),
	}
	s.metricCh <- e
	return nil
}

func (s *TelemetryServiceMongoAsyncImpl) Upload(ctx context.Context) error {
	for {
		batch := common.Batch(s.metricCh, s.batchInsertSize)
		if len(batch) == 0 {
			break
		}
		docs := make([]any, len(batch))
		for i, u := range batch {
			docs[i] = u
		}
		_, err := s.metricsCol.InsertMany(ctx, docs)
		if err != nil {
			return err
		}
	}
	for {
		batch := common.Batch(s.eventsCh, s.batchInsertSize)
		if len(batch) == 0 {
			break
		}
		docs := make([]any, len(batch))
		for i, u := range batch {
			docs[i] = u
		}
		_, err := s.metricsCol.InsertMany(ctx, docs)
		if err != nil {
			return err
		}
	}

	return nil
}
