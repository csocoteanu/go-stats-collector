package usecases

import (
	"context"
	"encoding/json"
	"log"
	"runtime"
	"time"

	"github.com/csocoteanu/go-stats-collector/stats_writer/domain"
	"github.com/shirou/gopsutil/mem"
)

type sensor struct {
	sensorID       string
	cassRepo       domain.CassandraRepository
	producer       domain.StatsProducer
	quit           chan struct{}
	reportInterval time.Duration
}

func NewSensor(
	sensorID string,
	reportInterval time.Duration,
	cassRepo domain.CassandraRepository,
	producer domain.StatsProducer) *sensor {

	s := sensor{
		sensorID:       sensorID,
		reportInterval: reportInterval,
		cassRepo:       cassRepo,
		producer:       producer,
		quit:           make(chan struct{}),
	}
	return &s
}

func (s sensor) Stop() {
	s.quit <- struct{}{}
}

func (s sensor) ScrapeStats() {
	ticker := time.NewTicker(s.reportInterval)
	ctx := context.Background()

	go func() {
		for {
			select {
			case <-ticker.C:
				stats, err := s.getSensorStats()
				if err != nil {
					panic(err)
				}

				statBytes, err := json.Marshal(stats)
				if err != nil {
					panic(err)
				}

				log.Printf("Polled: %s", string(statBytes))

				err = s.cassRepo.InsertHostStats(ctx, s.sensorID, time.Now(), &stats)
				if err != nil {
					panic(err)
				}

				if s.producer != nil {
					s.producer.Produce(stats)
				}
			case <-s.quit:
				log.Printf("Stopping scraper...")
				return
			}
		}
	}()
}

func (s sensor) getSensorStats() (domain.HostStats, error) {
	stats := domain.HostStats{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		ReportedTime: time.Now(),
	}

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return stats, err
	}

	stats.FreeMemoryMB = int(vmStat.Free / (1024 * 1024))
	stats.TotalMemoryMB = int(vmStat.Total / (1024 * 1024))

	return stats, nil
}
