package gateways

import (
	"context"
	"fmt"
	"time"

	"github.com/csocoteanu/go-stats-collector/stats_writer/domain"

	"github.com/gocql/gocql"
)

const (
	daySeconds                 = 24 * 60 * 60
	keyspace                   = "stats"
	liveSensorUpdatesTableName = "live_sensor_updates"
)

var insertLiveSensorUpdatesCql = fmt.Sprintf(`INSERT INTO %s (sensor_id, day, reported_time, os, arch, total_mem_mb, free_mem_mb) VALUES (?,?,?,?,?,?,?)`, liveSensorUpdatesTableName)

type CassandraStatsRepositoryConfig struct {
	ClusterIPs      []string
	Timeout         time.Duration
	NumConnections  int
	NumQueryRetries int
}

type CassandraStatsRepository struct {
	repoConfig    *CassandraStatsRepositoryConfig
	clusterConfig *gocql.ClusterConfig
	session       *gocql.Session
}

func NewCassandraStatsRepository(config *CassandraStatsRepositoryConfig) (*CassandraStatsRepository, error) {
	clusterConfig := gocql.NewCluster(config.ClusterIPs...)
	clusterConfig.Keyspace = keyspace
	clusterConfig.NumConns = config.NumConnections
	clusterConfig.Consistency = gocql.Quorum
	clusterConfig.MaxPreparedStmts = 100000
	clusterConfig.MaxRoutingKeyInfo = 100000
	clusterConfig.CQLVersion = "3.2.0"
	clusterConfig.Timeout = config.Timeout
	clusterConfig.ConnectTimeout = config.Timeout

	clusterConfig.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: config.NumQueryRetries}
	clusterConfig.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())

	session, err := clusterConfig.CreateSession()
	if err != nil {
		return nil, err
	}

	r := CassandraStatsRepository{
		repoConfig:    config,
		clusterConfig: clusterConfig,
		session:       session,
	}

	return &r, nil
}

func (r *CassandraStatsRepository) InsertHostStats(
	ctx context.Context,
	sensorID string,
	reportedTime time.Time,
	stats ...*domain.HostStats) error {

	day := int(reportedTime.Unix() / daySeconds)

	switch len(stats) {
	case 0:
		return nil
	case 1:
		s := stats[0]
		return r.session.
			Query(insertLiveSensorUpdatesCql,
				sensorID, day,
				s.ReportedTime,
				s.OS,
				s.Architecture,
				s.TotalMemoryMB,
				s.FreeMemoryMB).
			WithContext(ctx).
			Exec()
	default:
		batch := r.session.NewBatch(gocql.UnloggedBatch).WithContext(ctx)
		for _, s := range stats {
			batch.Query(insertLiveSensorUpdatesCql,
				sensorID, day,
				s.ReportedTime,
				s.OS,
				s.Architecture,
				s.TotalMemoryMB,
				s.FreeMemoryMB)
		}

		return r.session.ExecuteBatch(batch)
	}
}
