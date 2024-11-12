package collector

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryCollector struct{}

func CreateMemoryCollector() *MemoryCollector {
	return &MemoryCollector{}
}

func (c *MemoryCollector) Collect() ([]Metric, error) {
	var metrics []Metric
	now := time.Now()

	// Get virtual memory stats
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("error getting virtual memory metrics: %v", err)
	}

	// Create metrics
	metrics = append(metrics,
		// Overall memory usage percentage
		Metric{
			Name:      "memory_usage",
			Value:     vmem.UsedPercent,
			Timestamp: now,
			Labels:    map[string]string{"type": "virtual"},
		},
		// Total memory in bytes
		Metric{
			Name:      "memory_total",
			Value:     float64(vmem.Total),
			Timestamp: now,
			Labels:    map[string]string{"type": "virtual"},
		},
		// Used memory in bytes
		Metric{
			Name:      "memory_used",
			Value:     float64(vmem.Used),
			Timestamp: now,
			Labels:    map[string]string{"type": "virtual"},
		},
		// Free memory in bytes
		Metric{
			Name:      "memory_free",
			Value:     float64(vmem.Free),
			Timestamp: now,
			Labels:    map[string]string{"type": "virtual"},
		},
	)

	// Get swap memory stats
	swap, err := mem.SwapMemory()
	if err == nil { // Don't fail if swap metrics are unavailable
		metrics = append(metrics,
			// Swap usage percentage
			Metric{
				Name:      "memory_usage",
				Value:     swap.UsedPercent,
				Timestamp: now,
				Labels:    map[string]string{"type": "swap"},
			},
			// Total swap in bytes
			Metric{
				Name:      "memory_total",
				Value:     float64(swap.Total),
				Timestamp: now,
				Labels:    map[string]string{"type": "swap"},
			},
			// Used swap in bytes
			Metric{
				Name:      "memory_used",
				Value:     float64(swap.Used),
				Timestamp: now,
				Labels:    map[string]string{"type": "swap"},
			},
		)
	}

	// Add detailed virtual memory metrics if available
	if vmem.Cached > 0 {
		metrics = append(metrics,
			Metric{
				Name:      "memory_cached",
				Value:     float64(vmem.Cached),
				Timestamp: now,
				Labels:    map[string]string{"type": "virtual"},
			},
		)
	}

	if vmem.Buffers > 0 {
		metrics = append(metrics,
			Metric{
				Name:      "memory_buffers",
				Value:     float64(vmem.Buffers),
				Timestamp: now,
				Labels:    map[string]string{"type": "virtual"},
			},
		)
	}

	if vmem.Available > 0 {
		metrics = append(metrics,
			Metric{
				Name:      "memory_available",
				Value:     float64(vmem.Available),
				Timestamp: now,
				Labels:    map[string]string{"type": "virtual"},
			},
		)
	}

	return metrics, nil
}
