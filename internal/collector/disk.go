package collector

import (
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/disk"
)

type DiskCollector struct {
	lastStats   map[string]disk.IOCountersStat
	lastCollect time.Time
	mutex       sync.Mutex
}

func CreateDiskCollector() *DiskCollector {
	return &DiskCollector{
		lastStats: make(map[string]disk.IOCountersStat),
	}
}

func (c *DiskCollector) Collect() ([]Metric, error) {
	metrics := []Metric{}
	now := time.Now()

	mainDiskMountPoint := "/"

	usage, err := disk.Usage(mainDiskMountPoint)
	if err != nil {
		return nil, fmt.Errorf("error getting disk usage for main disk: %v", err)
	}

	metrics = append(metrics,
		Metric{
			Name:      "disk_usage",
			Value:     usage.UsedPercent,
			Timestamp: now,
			Labels: map[string]string{
				"mountpoint": mainDiskMountPoint,
				"fstype":     usage.Fstype,
			},
		},
		Metric{
			Name:      "disk_total",
			Value:     float64(usage.Total),
			Timestamp: now,
			Labels: map[string]string{
				"mountpoint": mainDiskMountPoint,
				"fstype":     usage.Fstype,
			},
		},
		Metric{
			Name:      "disk_free",
			Value:     float64(usage.Free),
			Timestamp: now,
			Labels: map[string]string{
				"mountpoint": mainDiskMountPoint,
				"fstype":     usage.Fstype,
			},
		},
	)

	ioStats, err := disk.IOCounters()
	if err != nil {
		fmt.Printf("Warning: error getting IO statistics: %v\n", err)
	} else {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		timeSinceLastCollect := now.Sub(c.lastCollect).Seconds()

		for device, stats := range ioStats {
			ioLabels := map[string]string{"device": device}

			if lastStat, exists := c.lastStats[device]; exists && timeSinceLastCollect > 0 {
				readSpeed := float64(stats.ReadBytes-lastStat.ReadBytes) / timeSinceLastCollect
				metrics = append(metrics, Metric{
					Name:      "disk_read_speed_bytes_per_second",
					Value:     readSpeed,
					Timestamp: now,
					Labels:    ioLabels,
				})

				writeSpeed := float64(stats.WriteBytes-lastStat.WriteBytes) / timeSinceLastCollect
				metrics = append(metrics, Metric{
					Name:      "disk_write_speed_bytes_per_second",
					Value:     writeSpeed,
					Timestamp: now,
					Labels:    ioLabels,
				})

				readIOPS := float64(stats.ReadCount-lastStat.ReadCount) / timeSinceLastCollect
				writeIOPS := float64(stats.WriteCount-lastStat.WriteCount) / timeSinceLastCollect

				metrics = append(metrics, []Metric{
					{
						Name:      "disk_read_iops",
						Value:     readIOPS,
						Timestamp: now,
						Labels:    ioLabels,
					},
					{
						Name:      "disk_write_iops",
						Value:     writeIOPS,
						Timestamp: now,
						Labels:    ioLabels,
					},
					{
						Name:      "disk_total_iops",
						Value:     readIOPS + writeIOPS,
						Timestamp: now,
						Labels:    ioLabels,
					},
				}...)
			}

			metrics = append(metrics, []Metric{
				{
					Name:      "disk_reads_total",
					Value:     float64(stats.ReadCount),
					Timestamp: now,
					Labels:    ioLabels,
				},
				{
					Name:      "disk_writes_total",
					Value:     float64(stats.WriteCount),
					Timestamp: now,
					Labels:    ioLabels,
				},
				{
					Name:      "disk_read_bytes_total",
					Value:     float64(stats.ReadBytes),
					Timestamp: now,
					Labels:    ioLabels,
				},
				{
					Name:      "disk_write_bytes_total",
					Value:     float64(stats.WriteBytes),
					Timestamp: now,
					Labels:    ioLabels,
				},
				{
					Name:      "disk_io_in_progress",
					Value:     float64(stats.IopsInProgress),
					Timestamp: now,
					Labels:    ioLabels,
				},
			}...)
		}

		c.lastStats = ioStats
		c.lastCollect = now
	}

	return metrics, nil
}
