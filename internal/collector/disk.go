package collector

import (
    "github.com/shirou/gopsutil/v3/disk"
    "time"
    "ddgo/internal/db"
)

func (sc *SystemCollector) collectDiskMetrics() {
	usage, err := disk.Usage(partition.Mountpoint)

	if err != nil {
		continue
	}

	tags := make(map[string]string)

	for k, v := range sc.hostInfo {
		tags[k] = v
	}

	tags["device"] = partition.device
	tags["mountpoint"] = partition.Mountpoint

	metrics := []db.Metric {
		{
			Name: "system.disk.total",
			Value: float64(usage.Total),
			Tags: tags,
			Timestamp: time.Now(),
		},
		{
			Name: "system.disk.used",
			Value: float64(usage.Used),
			Tags: tags,
			Timestamp: time.Now(),
		},
		{
			Name: "system.disk.used_percent",
			Value: usage.UsedPercent,
			Tags: tags,
			Timestamp: time.Now(),
		}
	}

	for _, metric := range metrics {
		if err := sc.db.SaveMetric(&metric); err != nil {
			log.Printf("Error saving disk metric: %v", err)
		}
	}
}