package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/csocoteanu/go-stats-collector/stats_writer/usecases"

	"github.com/csocoteanu/go-stats-collector/stats_writer/domain"

	"github.com/csocoteanu/go-stats-collector/stats_writer/gateways"
)

const defaultSensorCount = 10

var sensorCount = defaultSensorCount

func parseArgs() {
	flag.IntVar(&sensorCount, "sensor-count", defaultSensorCount, "sensors to start reporting data")
	flag.Parse()
}

func main() {
	parseArgs()

	config := gateways.CassandraStatsRepositoryConfig{
		ClusterIPs:      []string{"localhost"},
		Timeout:         30 * time.Second,
		NumQueryRetries: 3,
		NumConnections:  20,
	}

	cassRepo, err := gateways.NewCassandraStatsRepository(&config)
	if err != nil {
		panic(err)
	}
	log.Print("Created Cassandra repository...")

	timeStart := time.Now()

	sensors := []domain.Sensor{}
	for i := 0; i < sensorCount; i++ {
		sensor := usecases.NewSensor(fmt.Sprintf("%d", i), 500*time.Millisecond, cassRepo, nil)
		sensor.ScrapeStats()
		sensors = append(sensors, sensor)
	}
	log.Printf("Created %d sensors...", sensorCount)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case <-c:
		log.Printf("Quiting...")
		for _, sensor := range sensors {
			sensor.Stop()
		}
	}

	timeEnd := time.Now().Sub(timeStart)
	log.Printf("Running %d sensors for %s...", sensorCount, timeEnd.String())
}
