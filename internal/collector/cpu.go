package collector

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/process"
)

type CPUCollector struct {
	history     []Metric
	historySize int
}

func CreateCPUCollector(historySize int) *CPUCollector {
	return &CPUCollector{
		history:     make([]Metric, 0, historySize),
		historySize: historySize,
	}
}

func (c *CPUCollector) Collect() ([]Metric, error) {
	metrics := []Metric{}
	now := time.Now()

	// cpu usage
	usageMetrics, err := c.collectCPUUsage(now)
	if err != nil {
		return nil, err
	}
	metrics = append(metrics, usageMetrics...)

	// load avgs
	loadMetrics, err := c.collectLoadAverages(now)
	if err != nil {
		return nil, err
	}
	metrics = append(metrics, loadMetrics...)

	// cpu counts
	countMetrics, err := c.collectCPUCounts(now)
	if err != nil {
		return nil, err
	}
	metrics = append(metrics, countMetrics...)

	// cpu times
	timeMetrics, err := c.collectCPUTimes(now)
	if err != nil {
		return nil, err
	}
	metrics = append(metrics, timeMetrics...)

	// system stats
	systemMetrics, err := c.collectSystemStats(now)
	if err != nil {
		return nil, err
	}
	metrics = append(metrics, systemMetrics...)

	// store metrics in history
	// TODO: trend analysis
	c.storeInHistory(metrics)

	return metrics, nil
}

func (c *CPUCollector) collectCPUUsage(now time.Time) ([]Metric, error) {
	percentages, err := cpu.Percent(time.Second, true)
	if err != nil {
		return nil, fmt.Errorf("error collecting CPU usage: %v", err)
	}

	metrics := []Metric{}
	for i, usage := range percentages {
		metrics = append(metrics, Metric{
			Name:      "cpu_usage",
			Value:     usage,
			Timestamp: now,
			Labels:    map[string]string{"cpu": fmt.Sprintf("cpu%d", i)},
		})
	}

	return metrics, nil
}

func (c *CPUCollector) collectLoadAverages(now time.Time) ([]Metric, error) {
	loadAvg, err := load.Avg()
	if err != nil {
		return nil, fmt.Errorf("error collecting load average: %v", err)
	}

	metrics := []Metric{
		{
			Name:      "cpu_load_average_1m",
			Value:     loadAvg.Load1,
			Timestamp: now,
			Labels:    map[string]string{"interval": "1m"},
		},
		{
			Name:      "cpu_load_average_5m",
			Value:     loadAvg.Load5,
			Timestamp: now,
			Labels:    map[string]string{"interval": "5m"},
		},
		{
			Name:      "cpu_load_average_15m",
			Value:     loadAvg.Load15,
			Timestamp: now,
			Labels:    map[string]string{"interval": "15m"},
		},
	}

	return metrics, nil
}

func (c *CPUCollector) collectCPUCounts(now time.Time) ([]Metric, error) {
	counts, err := cpu.Counts(true)
	if err != nil {
		return nil, fmt.Errorf("error collecting CPU counts: %v", err)
	}

	physicalCounts, err := cpu.Counts(false)
	if err != nil {
		return nil, fmt.Errorf("error collecting physical CPU counts: %v", err)
	}

	metrics := []Metric{
		{
			Name:      "cpu_cores_logical",
			Value:     float64(counts),
			Timestamp: now,
			Labels:    map[string]string{},
		},
		{
			Name:      "cpu_cores_physical",
			Value:     float64(physicalCounts),
			Timestamp: now,
			Labels:    map[string]string{},
		},
	}

	if physicalCounts > 0 {
		metrics = append(metrics, Metric{
			Name:      "cpu_hyperthread_ratio",
			Value:     float64(counts) / float64(physicalCounts),
			Timestamp: now,
			Labels:    map[string]string{},
		})
	}

	return metrics, nil
}

func (c *CPUCollector) collectCPUTimes(now time.Time) ([]Metric, error) {
	times, err := cpu.Times(true)
	if err != nil {
		return nil, fmt.Errorf("error collecting CPU times: %v", err)
	}

	metrics := []Metric{}
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

	return metrics, nil
}

func getSystemStats() (uint64, uint64, error) {
	_, _, _, err := host.PlatformInformation()
	if err != nil {
		return 0, 0, err
	}

	cpuStats, err := cpu.Times(false)
	if err != nil {
		return 0, 0, err
	}

	// Note: The actual values might need adjustment based on your OS
	// These are approximations based on CPU statistics
	contextSwitches := uint64(cpuStats[0].Irq + cpuStats[0].Softirq)
	interrupts := uint64(cpuStats[0].Irq)

	return contextSwitches, interrupts, nil
}

func (c *CPUCollector) collectSystemStats(now time.Time) ([]Metric, error) {
	contextSwitches, interrupts, err := getSystemStats()
	if err != nil {
		return nil, fmt.Errorf("error collecting context switches: %v", err)
	}

	metrics := []Metric{
		{
			Name:      "cpu_context_switches_total",
			Value:     float64(contextSwitches),
			Timestamp: now,
			Labels:    map[string]string{},
		},
		{
			Name:      "cpu_interrupts_total",
			Value:     float64(interrupts),
			Timestamp: now,
			Labels:    map[string]string{},
		},
	}

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

	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("error collecting process stats: %v", err)
	}

	metrics = append(metrics, Metric{
		Name:      "system_processes_total",
		Value:     float64(len(processes)),
		Timestamp: now,
		Labels:    map[string]string{},
	})

	var (
		totalThreads int32
		running      int
		sleeping     int
		stopped      int
		zombie       int
	)

	for _, p := range processes {
		if numThreads, err := p.NumThreads(); err == nil {
			totalThreads += numThreads
		}

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

func (c *CPUCollector) storeInHistory(metrics []Metric) {
	c.history = append(c.history, metrics...)
	if len(c.history) > c.historySize {
		c.history = c.history[len(c.history)-c.historySize:]
	}
}
