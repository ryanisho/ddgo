// internal/collector/disk.go
package collector

import (
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

type DiskCollector struct {
	lastIOCounters map[string]disk.IOCountersStat
	lastCheck      time.Time
}

func NewDiskCollector() *DiskCollector {
	return &DiskCollector{
		lastIOCounters: make(map[string]disk.IOCountersStat),
		lastCheck:      time.Now(),
	}
}

func (c *DiskCollector) Collect() ([]Metric, error) {
	var metrics []Metric
	now := time.Now()

	// Get all partitions
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("error getting disk partitions: %v", err)
	}

	// Collect metrics for each partition
	for _, partition := range partitions {
		// Skip special filesystems
		if shouldSkipFilesystem(partition.Fstype) {
			continue
		}

		// Get usage statistics
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue // Skip this partition if we can't get usage
		}

		// Add disk space metrics
		metrics = append(metrics,
			// Usage percentage
			Metric{
				Name:      "disk_usage",
				Value:     usage.UsedPercent,
				Timestamp: now,
				Labels: map[string]string{
					"device":     partition.Device,
					"mountpoint": partition.Mountpoint,
					"fstype":     partition.Fstype,
				},
			},
			// Total space
			Metric{
				Name:      "disk_total",
				Value:     float64(usage.Total),
				Timestamp: now,
				Labels: map[string]string{
					"device":     partition.Device,
					"mountpoint": partition.Mountpoint,
					"fstype":     partition.Fstype,
				},
			},
			// Free space
			Metric{
				Name:      "disk_free",
				Value:     float64(usage.Free),
				Timestamp: now,
				Labels: map[string]string{
					"device":     partition.Device,
					"mountpoint": partition.Mountpoint,
					"fstype":     partition.Fstype,
				},
			},
		)

		// Add inode metrics if available (Unix/Linux)
		if usage.InodesTotal > 0 {
			metrics = append(metrics,
				Metric{
					Name:      "disk_inodes_usage",
					Value:     usage.InodesUsedPercent,
					Timestamp: now,
					Labels: map[string]string{
						"device":     partition.Device,
						"mountpoint": partition.Mountpoint,
						"fstype":     partition.Fstype,
					},
				},
			)
		}
	}

	// Get IO statistics
	ioCounters, err := disk.IOCounters()
	if err == nil { // Don't fail if IO metrics are unavailable
		timeDiff := now.Sub(c.lastCheck).Seconds()

		for deviceName, ioStats := range ioCounters {
			// Calculate rates if we have previous measurements
			if lastStats, exists := c.lastIOCounters[deviceName]; exists && timeDiff > 0 {
				// Read rate in bytes per second
				readRate := float64(ioStats.ReadBytes-lastStats.ReadBytes) / timeDiff
				// Write rate in bytes per second
				writeRate := float64(ioStats.WriteBytes-lastStats.WriteBytes) / timeDiff
				// IO operations per second
				iopsRead := float64(ioStats.ReadCount-lastStats.ReadCount) / timeDiff
				iopsWrite := float64(ioStats.WriteCount-lastStats.WriteCount) / timeDiff

				metrics = append(metrics,
					// IO rates
					Metric{
						Name:      "disk_read_rate",
						Value:     readRate,
						Timestamp: now,
						Labels: map[string]string{
							"device": deviceName,
							"unit":   "bytes/s",
						},
					},
					Metric{
						Name:      "disk_write_rate",
						Value:     writeRate,
						Timestamp: now,
						Labels: map[string]string{
							"device": deviceName,
							"unit":   "bytes/s",
						},
					},
					// IOPS
					Metric{
						Name:      "disk_iops_read",
						Value:     iopsRead,
						Timestamp: now,
						Labels: map[string]string{
							"device": deviceName,
						},
					},
					Metric{
						Name:      "disk_iops_write",
						Value:     iopsWrite,
						Timestamp: now,
						Labels: map[string]string{
							"device": deviceName,
						},
					},
				)
			}
		}
		// Store current counters for next collection
		c.lastIOCounters = ioCounters
	}

	c.lastCheck = now
	return metrics, nil
}

// shouldSkipFilesystem returns true for special filesystems that should be ignored
func shouldSkipFilesystem(fstype string) bool {
	// List of filesystem types to skip
	skipFs := map[string]bool{
		"devfs":      true,
		"tmpfs":      true,
		"devtmpfs":   true,
		"squashfs":   true,
		"iso9660":    true,
		"overlay":    true,
		"aufs":       true,
		"proc":       true,
		"sysfs":      true,
		"devpts":     true,
		"securityfs": true,
		"cgroup":     true,
		"cgroup2":    true,
		"pstore":     true,
		"debugfs":    true,
		"hugetlbfs":  true,
		"mqueue":     true,
		"fusectl":    true,
	}

	return skipFs[strings.ToLower(fstype)]
}

// GetDiskAlerts returns warnings about disk usage
func (c *DiskCollector) GetDiskAlerts() []string {
	var alerts []string

	partitions, err := disk.Partitions(false)
	if err != nil {
		return alerts
	}

	for _, partition := range partitions {
		if shouldSkipFilesystem(partition.Fstype) {
			continue
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		// Alert on high disk usage
		if usage.UsedPercent > 90 {
			alerts = append(alerts, fmt.Sprintf(
				"Critical: High disk usage on %s: %.1f%% used",
				partition.Mountpoint,
				usage.UsedPercent,
			))
		} else if usage.UsedPercent > 80 {
			alerts = append(alerts, fmt.Sprintf(
				"Warning: Disk usage on %s: %.1f%% used",
				partition.Mountpoint,
				usage.UsedPercent,
			))
		}

		// Alert on inode usage (Unix/Linux)
		if usage.InodesTotal > 0 && usage.InodesUsedPercent > 90 {
			alerts = append(alerts, fmt.Sprintf(
				"Warning: High inode usage on %s: %.1f%% used",
				partition.Mountpoint,
				usage.InodesUsedPercent,
			))
		}
	}

	return alerts
}
