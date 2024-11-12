package collector

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/process"
)

type CPUCollector struct{}

func CreateCPUCollector() *CPUCollector {
	return &CPUCollector{}
}

func getSystemStats() (uint64, uint64, error) {
	_, _, _, err := host.PlatformInformation()
	if err != nil {
		return 0, 0, err
	}

	// Get context switches and interrupts from CPU stats
	cpuStats, err := cpu.Times(false) // false for total stats
	if err != nil {
		return 0, 0, err
	}

	// Note: The actual values might need adjustment based on your OS
	// These are approximations based on CPU statistics
	contextSwitches := uint64(cpuStats[0].Irq + cpuStats[0].Softirq)
	interrupts := uint64(cpuStats[0].Irq)

	return contextSwitches, interrupts, nil
}

func (c *CPUCollector) Collect() ([]Metric, error) {
	metrics := []Metric{}
	now := time.Now()

	// 1. CPU Usage
	percentages, err := cpu.Percent(time.Second, true)
	if err != nil {
		return nil, fmt.Errorf("error collecting CPU usage: %v", err)
	}

	for i, usage := range percentages {
		metrics = append(metrics, Metric{
			Name:      "cpu_usage",
			Value:     usage,
			Timestamp: now,
			Labels:    map[string]string{"cpu": fmt.Sprintf("cpu%d", i)},
		})
	}

	// 2. CPU Load Average
	loadAvg, err := load.Avg()
	if err != nil {
		return nil, fmt.Errorf("error collecting load average: %v", err)
	}

	metrics = append(metrics, Metric{
		Name:      "cpu_load_average_1m",
		Value:     loadAvg.Load1,
		Timestamp: now,
		Labels:    map[string]string{"interval": "1m"},
	})

	metrics = append(metrics, Metric{
		Name:      "cpu_load_average_5m",
		Value:     loadAvg.Load5,
		Timestamp: now,
		Labels:    map[string]string{"interval": "5m"},
	})

	metrics = append(metrics, Metric{
		Name:      "cpu_load_average_15m",
		Value:     loadAvg.Load15,
		Timestamp: now,
		Labels:    map[string]string{"interval": "15m"},
	})

	// log.Printf("Metrics: %v", metrics)

	// 3. CPU Frequency and Info
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("error collecting CPU frequency: %v", err)
	}

	for i, info := range cpuInfo {
		metrics = append(metrics, Metric{
			Name:      "cpu_frequency_current",
			Value:     float64(info.Mhz),
			Timestamp: now,
			Labels: map[string]string{
				"cpu":    fmt.Sprintf("cpu%d", i),
				"model":  info.ModelName,
				"vendor": info.VendorID,
			},
		})

		metrics = append(metrics, Metric{
			Name:      "cpu_cores",
			Value:     float64(info.Cores),
			Timestamp: now,
			Labels: map[string]string{
				"cpu":    fmt.Sprintf("cpu%d", i),
				"model":  info.ModelName,
				"vendor": info.VendorID,
			},
		})

		if info.CacheSize > 0 {
			metrics = append(metrics, Metric{
				Name:      "cpu_cache_size_mb",
				Value:     float64(info.CacheSize) / 1024.0,
				Timestamp: now,
				Labels: map[string]string{
					"cpu":    fmt.Sprintf("cpu%d", i),
					"model":  info.ModelName,
					"vendor": info.VendorID,
				},
			})
		}
	}

	// 4. CPU States (Time spent in different modes)
	times, err := cpu.Times(true) // true for per-CPU stats
	if err != nil {
		return nil, fmt.Errorf("error collecting CPU times: %v", err)
	}

	for i, cpuTime := range times {
		cpuLabels := map[string]string{"cpu": fmt.Sprintf("cpu%d", i)}

		metrics = append(metrics, Metric{
			Name:      "cpu_time_user",
			Value:     cpuTime.User,
			Timestamp: now,
			Labels:    cpuLabels,
		})

		metrics = append(metrics, Metric{
			Name:      "cpu_time_system",
			Value:     cpuTime.System,
			Timestamp: now,
			Labels:    cpuLabels,
		})

		metrics = append(metrics, Metric{
			Name:      "cpu_time_idle",
			Value:     cpuTime.Idle,
			Timestamp: now,
			Labels:    cpuLabels,
		})

		metrics = append(metrics, Metric{
			Name:      "cpu_time_iowait",
			Value:     cpuTime.Iowait,
			Timestamp: now,
			Labels:    cpuLabels,
		})

		metrics = append(metrics, Metric{
			Name:      "cpu_time_irq",
			Value:     cpuTime.Irq + cpuTime.Softirq,
			Timestamp: now,
			Labels:    cpuLabels,
		})
	}

	// 5. Context Switches and Interrupts
	contextSwitches, interrupts, err := getSystemStats()
	if err != nil {
		return nil, fmt.Errorf("error collecting context switches: %v", err)
	}

	metrics = append(metrics, Metric{
		Name:      "cpu_context_switches_total",
		Value:     float64(contextSwitches),
		Timestamp: now,
		Labels:    map[string]string{},
	})

	metrics = append(metrics, Metric{
		Name:      "cpu_interrupts_total",
		Value:     float64(interrupts),
		Timestamp: now,
		Labels:    map[string]string{},
	})

	// System boot time
	bootTime, err := host.BootTime()
	if err != nil {
		return nil, fmt.Errorf("error collecting boot time: %v", err)
	}

	metrics = append(metrics, Metric{
		Name:      "system_boot_time_seconds",
		Value:     float64(bootTime),
		Timestamp: now,
		Labels:    map[string]string{},
	})

	// 6. Process and Thread Statistics
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("error collecting process stats: %v", err)
	}

	// Process count
	metrics = append(metrics, Metric{
		Name:      "system_processes_total",
		Value:     float64(len(processes)),
		Timestamp: now,
		Labels:    map[string]string{},
	})

	// Thread count and process states
	var (
		totalThreads int32
		running      int
		sleeping     int
		stopped      int
		zombie       int
	)

	for _, p := range processes {
		// Count threads
		if numThreads, err := p.NumThreads(); err == nil {
			totalThreads += numThreads
		}

		// Count process states
		if status, err := p.Status(); err == nil {
			switch status[0] {
			case 'R':
				running++
			case 'S':
				sleeping++
			case 'T':
				stopped++
			case 'Z':
				zombie++
			}
		}
	}

	metrics = append(metrics, Metric{
		Name:      "system_threads_total",
		Value:     float64(totalThreads),
		Timestamp: now,
		Labels:    map[string]string{},
	})

	// Process states
	metrics = append(metrics, Metric{
		Name:      "system_processes_state",
		Value:     float64(running),
		Timestamp: now,
		Labels:    map[string]string{"state": "running"},
	})

	metrics = append(metrics, Metric{
		Name:      "system_processes_state",
		Value:     float64(sleeping),
		Timestamp: now,
		Labels:    map[string]string{"state": "sleeping"},
	})

	metrics = append(metrics, Metric{
		Name:      "system_processes_state",
		Value:     float64(stopped),
		Timestamp: now,
		Labels:    map[string]string{"state": "stopped"},
	})

	metrics = append(metrics, Metric{
		Name:      "system_processes_state",
		Value:     float64(zombie),
		Timestamp: now,
		Labels:    map[string]string{"state": "zombie"},
	})

	return metrics, nil
}
