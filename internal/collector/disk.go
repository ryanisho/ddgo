package collector

import (
	"time"

	"github.com/shirou/gopsutil/disk" // disk package
)

type DiskCollector struct{}

func (c *DiskCollector) Collect() ([]Metric, error) {
	usage, err := disk.Usage("/")

	if err != nil {
		return nil, err
	}

	now := time.Now()

	return []Metric{
		{
			Name:      "disk_usage",
			Value:     usage.UsedPercent,
			Timestamp: now,
			Labels:    map[string]string{"path": "/"},
		},
	}, nil
}
