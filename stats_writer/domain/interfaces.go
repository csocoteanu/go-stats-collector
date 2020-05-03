package domain

import (
	"context"
	"time"
)

//go:generate mockgen -destination=./mocks/mock_domain_interfaces.go -package=test -source=interfaces.go

// CassandraRepository ...
type CassandraRepository interface {
	// InsertHostStats inserts the provided stats associated for the given sensor ID
	InsertHostStats(ctx context.Context, sensorID string, reportedTime time.Time, stats ...*HostStats) error
}

// Sensor is the sensor which is being scrapped for data
type Sensor interface {
	// Stop stops the sensor
	Stop()
	// Starts the sensor scrapping stats
	ScrapeStats()
}

// StatsProducer ...
type StatsProducer interface {
	// Produce produces a message
	Produce(msg HostStats)
}

// StatsConsumer ...
type StatsConsumer interface {
	// Consume consumes a message
	Consume() (HostStats, bool)
}
