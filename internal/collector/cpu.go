package collector 

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"time"
	"ddgo/internal/db"
)

func (sc *SystemCollector) collectCPUMetrics() {
	percentages, err := cpu.Percent(0, false)

	if err != nil {
		return
	}

	if len(percentages) > 0 {
		metric := &db.Metric{
			Name: "system.cpu.usage",
			Value: percentages[0],
			Tags: sc.hostInfo,
			Timestamp: time.Now(),
		}
	}

	if err := sc.db.SaveMetric(metric); err != nill {
		log.Printf("Error saving CPU metric: %v", err)
	}
}