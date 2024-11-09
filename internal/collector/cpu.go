package collector

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu" // cpu package
)

type CPUCollector struct{}

// func NewCPUCollector() *CPUCollector {
// 	return &CPUCollector{}
// }

func (c *CPUCollector) Collect() ([]Metric, error) {
	percentages, err := cpu.Percent(time.Second, true)

	if err != nil {
		return nil, err
	}

	metrics := make([]Metric, len(percentages))
	now := time.Now()

	for i, usage := range percentages {
		metrics[i] = Metric{
			Name:      "cpu_usage",
			Value:     usage,
			Timestamp: now,
			Labels:    map[string]string{"cpu": fmt.Sprintf("cpu%d", i)},
		}
	}

	return metrics, nil
}
