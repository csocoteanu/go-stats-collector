package domain

import "time"

type HostStats struct {
	ReportedTime  time.Time `json:"reported_time"`
	OS            string    `json:"os"`
	Architecture  string    `json:"architecture"`
	TotalMemoryMB int       `json:"total_memory_mb"`
	FreeMemoryMB  int       `json:"free_memory_mb"`
}
