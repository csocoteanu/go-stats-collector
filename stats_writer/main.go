package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/csocoteanu/go-stats-collector/stats_writer/gateways"

	"github.com/csocoteanu/go-stats-collector/stats_writer/domain"

	"github.com/shirou/gopsutil/mem"
)

func getHostStats() (domain.HostStats, error) {
	s := domain.HostStats{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		ReportedTime: time.Now(),
	}

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return s, err
	}

	s.FreeMemoryMB = int(vmStat.Free / (1024 * 1024))
	s.TotalMemoryMB = int(vmStat.Total / (1024 * 1024))

	return s, nil
}

func scrapeStats(quit chan struct{}, cassRepo *gateways.CassandraRepository) {
	ticker := time.NewTicker(500 * time.Millisecond)
	ctx := context.Background()

	for {
		select {
		case <-ticker.C:
			stats, err := getHostStats()
			if err != nil {
				panic(err)
			}

			statBytes, err := json.Marshal(stats)
			if err != nil {
				panic(err)
			}

			log.Printf("Polled: %s", string(statBytes))

			err = cassRepo.InsertHostStats(ctx, "0", time.Now(), &stats)
			if err != nil {
				panic(err)
			}
		case <-quit:
			log.Printf("Stopping scraper...")
			return
		}
	}
}

func main() {
	config := gateways.CassandraRepositoryConfig{
		ClusterIPs:      []string{"localhost"},
		Timeout:         30 * time.Second,
		NumQueryRetries: 3,
		NumConnections:  20,
	}

	cassRepo, err := gateways.NewCassandraRepository(&config)
	if err != nil {
		panic(err)
	}

	quit := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go scrapeStats(quit, cassRepo)

	select {
	case <-c:
		log.Printf("Quiting...")
		quit <- struct{}{}
	}
}
