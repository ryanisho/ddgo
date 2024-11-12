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

	vmem, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("error getting virtual memory metrics: %v", err)
	}

	metrics = append(metrics,
		Metric{
			Name:      "memory_usage",
			Value:     vmem.UsedPercent,
			Timestamp: now,
			Labels:    map[string]string{"type": "virtual"},
		},
		Metric{
			Name:      "memory_total",
			Value:     float64(vmem.Total),
			Timestamp: now,
			Labels:    map[string]string{"type": "virtual"},
		},
		Metric{
			Name:      "memory_used",
			Value:     float64(vmem.Used),
			Timestamp: now,
			Labels:    map[string]string{"type": "virtual"},
		},
		Metric{
			Name:      "memory_free",
			Value:     float64(vmem.Free),
			Timestamp: now,
			Labels:    map[string]string{"type": "virtual"},
		},
	)

	metrics = append(metrics, []Metric{
		{
			Name:      "memory_page_tables",
			Value:     float64(vmem.PageTables),
			Timestamp: now,
			Labels:    map[string]string{"type": "paged"},
		},
		{
			Name:      "memory_mapped",
			Value:     float64(vmem.Mapped),
			Timestamp: now,
			Labels:    map[string]string{"type": "paged"},
		},
		{
			Name:      "memory_slab",
			Value:     float64(vmem.Slab),
			Timestamp: now,
			Labels:    map[string]string{"type": "paged"},
		},
		{
			Name:      "memory_page_cache",
			Value:     float64(vmem.Cached),
			Timestamp: now,
			Labels:    map[string]string{"type": "paged"},
		},
		{
			Name:      "memory_writeback_temp",
			Value:     float64(vmem.WriteBackTmp),
			Timestamp: now,
			Labels:    map[string]string{"type": "paged"},
		},
		{
			Name:      "memory_dirty_pages",
			Value:     float64(vmem.Dirty),
			Timestamp: now,
			Labels:    map[string]string{"type": "paged"},
		},
		{
			Name:      "memory_writeback_pages",
			Value:     float64(vmem.WriteBack),
			Timestamp: now,
			Labels:    map[string]string{"type": "paged"},
		},
	}...)

	swap, err := mem.SwapMemory()
	if err == nil {
		metrics = append(metrics, []Metric{
			{
				Name:      "memory_usage",
				Value:     swap.UsedPercent,
				Timestamp: now,
				Labels:    map[string]string{"type": "swap"},
			},
			{
				Name:      "memory_total",
				Value:     float64(swap.Total),
				Timestamp: now,
				Labels:    map[string]string{"type": "swap"},
			},
			{
				Name:      "memory_used",
				Value:     float64(swap.Used),
				Timestamp: now,
				Labels:    map[string]string{"type": "swap"},
			},
			{
				Name:      "swap_in_bytes_total",
				Value:     float64(swap.Sin),
				Timestamp: now,
				Labels:    map[string]string{"type": "swap"},
			},
			{
				Name:      "swap_out_bytes_total",
				Value:     float64(swap.Sout),
				Timestamp: now,
				Labels:    map[string]string{"type": "swap"},
			},
			{
				Name:      "swap_free",
				Value:     float64(swap.Free),
				Timestamp: now,
				Labels:    map[string]string{"type": "swap"},
			},
		}...)
	}

	return metrics, nil
}
