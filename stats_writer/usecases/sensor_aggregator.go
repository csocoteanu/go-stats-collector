package usecases

import (
	"context"
	"time"

	"github.com/csocoteanu/go-stats-collector/stats_writer/domain"
)

type sensorAggregator struct {
	sensorID string
	consumer domain.StatsConsumer
	cassRepo domain.CassandraRepository
	capacity int
	aggStats []domain.HostStats
}

func NewSensorAggregator(
	sensorID string,
	consumer domain.StatsConsumer,
	capacity int,
	repository domain.CassandraRepository) *sensorAggregator {

	sa := sensorAggregator{
		sensorID: sensorID,
		consumer: consumer,
		capacity: capacity,
		cassRepo: repository,
		aggStats: make([]domain.HostStats, capacity),
	}

	sa.aggStats = sa.aggStats[:0]

	return &sa
}

func (sa sensorAggregator) ListenForEvents() {
	for {
		stat, hasStat := sa.consumer.Consume()
		if !hasStat {
			break
		}

		sa.aggStats = append(sa.aggStats, stat)
		if len(sa.aggStats) == sa.capacity {
			err := sa.aggregateStats()
			if err != nil {
				panic(err)
			}

			sa.aggStats = sa.aggStats[:0]
		}
	}

	if len(sa.aggStats) > 0 {
		err := sa.aggregateStats()
		if err != nil {
			panic(err)
		}
	}
}

func (sa sensorAggregator) aggregateStats() error {
	result := domain.AggregateStats(sa.aggStats...)
	err := sa.cassRepo.InsertHostStats(context.Background(), sa.sensorID, time.Now(), &result)
	return err
}
