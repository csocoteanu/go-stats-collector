package domain

import "time"

type HostStats struct {
	ReportedTime  time.Time `json:"reported_time"`
	OS            string    `json:"os"`
	Architecture  string    `json:"architecture"`
	TotalMemoryMB int       `json:"total_memory_mb"`
	FreeMemoryMB  int       `json:"free_memory_mb"`
}

func AggregateStats(stats ...HostStats) HostStats {
	if len(stats) == 0 {
		return HostStats{}
	}

	s := HostStats{
		OS:           stats[0].OS,
		Architecture: stats[0].Architecture,
		ReportedTime: stats[len(stats)/2].ReportedTime,
	}

	var totalMemMB int
	var freeMemMB int

	for _, s := range stats {
		totalMemMB += s.TotalMemoryMB
		freeMemMB += s.FreeMemoryMB
	}

	s.TotalMemoryMB = totalMemMB / len(stats)
	s.FreeMemoryMB = freeMemMB / len(stats)

	return s
}
