package collector

import (
    "github.com/shirou/gopsutil/v3/mem"
    "time"
    "ddgo/internal/db"
)

func (sc *SystemCollector) collectMemoryMetrics() {
	vmStat, err := mem.VirtualMemory()

	if err != nil {
		return
	}

	metrics := []db.Metric[] {
		{
			Name: "system.mem.total",
			Value: float64(vmStat.Total),
			Tags: sc.hostInfo,
			Timestamp: time.Now(),
		},
		{
			Name: "system.mem.used",
			Value: float64(vmStat.Used),
			Tags: sc.hostInfo,
			Timestamp: time.Now(),
		},
		{
			Name: "system.mem.used_percent",
			Value: vmStat.UsedPercent,
			Tags: sc.hostInfo,
			Timestamp: time.Now(),
		},
	}

	for _, metric := range metrics {
		if err := sc.db.SaveMetric(&metric) err != nil {
			log.Printf("Error saving memory metric: %v", err)
		}
	}
}


