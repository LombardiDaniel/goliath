package services

import (
	"context"
	"sync"
	"time"

	"github.com/LombardiDaniel/gopherbase/common"
	"github.com/LombardiDaniel/gopherbase/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TelemetryServiceMongoAsyncImpl struct {
	mongoClient     *mongo.Client
	metricsCol      *mongo.Collection
	eventsCol       *mongo.Collection
	batchInsertSize uint32
	metricCh        chan models.Metric
	eventsCh        chan models.Event
	counters        []Counter
}

type CounterMongoAsyncImpl struct {
	metricsCol *mongo.Collection
	metricName string
	tags       map[string]string
	val        uint64
	valLock    sync.Mutex
}

func (c *CounterMongoAsyncImpl) Increment(count uint64) {
	c.valLock.Lock()
	defer c.valLock.Unlock()

	c.val += count
}

func (c *CounterMongoAsyncImpl) Upload(ctx context.Context) error {
	filter := bson.M{
		"name": c.metricName,
		"tags": c.tags,
	}

	c.valLock.Lock()
	v := c.val
	c.val = 0
	c.valLock.Unlock()

	update := bson.M{"value": bson.M{"$inc": v}}
	upsert := true
	_, err := c.metricsCol.UpdateOne(ctx, filter, update, &options.UpdateOptions{Upsert: &upsert})
	return err
}

func NewTelemetryServiceMongoAsyncImpl(mongoClient *mongo.Client, metricsCol, eventsCol *mongo.Collection, batchInsertSize uint32) TelemetryService {
	return &TelemetryServiceMongoAsyncImpl{
		mongoClient:     mongoClient,
		metricsCol:      metricsCol,
		eventsCol:       eventsCol,
		metricCh:        make(chan models.Metric),
		eventsCh:        make(chan models.Event),
		batchInsertSize: batchInsertSize,
		counters:        []Counter{},
	}
}

func (s *TelemetryServiceMongoAsyncImpl) GetCounter(ctx context.Context, metricName string, tags map[string]string) (Counter, error) {
	c := &CounterMongoAsyncImpl{
		metricsCol: s.metricsCol,
		metricName: metricName,
		tags:       tags,
	}

	s.counters = append(s.counters, c)

	return c, nil
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
	for _, c := range s.counters {
		err := c.Upload(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
