package metrics

import (
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/mem"
    "time"
)

type SystemMetrics struct {
	CPUUsage float64
	MemoryUsage uint64
}

func GetSystemMetrics() SystemMetrics {
	cpuPercent, _ := cpu.Percent(0, false)
	memStats, _ := mem.VirtualMemory()

	return SystemMetrics {
		CPUUsage: cpuPercent[0],
		MemoryUsage: memStats.Used,
	}
}