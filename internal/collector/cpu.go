package collector

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/v3/cpu"
)

type cpuCollector struct {
	mu            sync.RWMutex // thread safety
	config        CPUMetricConfig
	prevTimes     []cpu.TimesStat
	prevTimestamp time.Time
	loadHistory   []float64
	stats         CPUStats
	alerts        []Alert
}

// of type CPUMetricConfig, inerface CPUCollector
func CreateParameterizedCPUCollector(config CPUMetricConfig) CPUCollector {
	return &cpuCollector{
		config:      config,
		loadHistory: make([]float64, 0, config.HistorySize),
	}
}

func CreateDefaultCPUCollector() CPUCollector {
	return CreateParameterizedCPUCollector(CPUMetricConfig{
		CollectTemperature: false,
		CollectFrequency:   true,
		CollectLoadAvg:     true,
		HistorySize:        DefaultHistorySize,
		UsageThreshold:     DefaultUsageThreshold,
		TempThreshold:      DefaultTempThreshold,
		LoadThreshold:      DefaultLoadThreshold,
	})
}

func (c *cpuCollector) Collect() ([]Metric, error) {
	var metrics []Metric
	now := time.Now()

	metricsChan := make(chan []Metric, 4)
	errorsChan := make(chan error, 4)

	// collect metrics concurrently
	go c.collectUsage(now, metricsChan, errorsChan)
	go c.collectFrequency(now, metricsChan, errorsChan)
	// collect temperature here
	go c.collectLoad(now, metricsChan, errorsChan)

	for i := 0; i < 3; i++ {
		select {
		case newMetrics := <-metricsChan:
			metrics = append(metrics, newMetrics...)
		case err := <-errorsChan:
			if err != nil {
				c.addAlert(Alert{
					Source:    MetricCPUUsage,
					Level:     AlertLevelWarning,
					Message:   fmt.Sprintf("Collection error: %v", err),
					Timestamp: now,
				})
			}
		}
	}

	return metrics, nil
}

func (c *cpuCollector) collectUsage(now time.Time, metrics chan<- []Metric, errors chan<- error) {
	var localMetrics []Metric

	// cpu usage per core
	perCPU, err := cpu.Percent(time.Second, true)
	if err != nil {
		errors <- fmt.Errorf("per-core CPU Usage: %v", err)
		return
	}

	// total cpu usage
	totalCPU, err := cpu.Percent(time.Second, false)
	if err != nil {
		errors <- fmt.Errorf("total CPU Usage: %v", err)
		return
	}

	// writeback
	c.mu.Lock()
	c.stats.Usage = totalCPU[0] // [0] = %
	c.mu.Unlock()

	// collect total cpu
	localMetrics = append(localMetrics, Metric{
		Name:      MetricCPUUsage,
		Value:     totalCPU[0],
		Timestamp: now,
		Labels:    map[string]string{"type": "total"},
	})

	// collect per core
	for i, usage := range perCPU {
		localMetrics = append(localMetrics, Metric{
			Name:      MetricCPUUsage,
			Value:     usage,
			Timestamp: now,
			Labels: map[string]string{
				"type": "core",
				"core": fmt.Sprintf("%d", i),
			},
		})

		if usage > c.config.UsageThreshold {
			c.addAlert(Alert{
				Source:    MetricCPUUsage,
				Level:     AlertLevelWarning,
				Message:   fmt.Sprintf("High CPU Usage on core %d: %.2f%%", i, usage),
				Timestamp: now,
			})
		}
	}

	metrics <- localMetrics
	errors <- nil
}

func (c *cpuCollector) collectFrequency(now time.Time, metrics chan<- []Metric, errors chan<- error) {
	if !c.config.CollectFrequency {
		metrics <- nil
		errors <- nil
		return
	}

	var localMetrics []Metric

	info, err := cpu.Info()
	if err != nil {
		errors <- fmt.Errorf("CPU frequency: %v", err)
		return
	}

	for i, cpu := range info {
		freq := float64(cpu.Mhz)

		localMetrics = append(localMetrics, Metric{
			Name:      MetricCPUFreq,
			Value:     freq,
			Timestamp: now,
			Labels: map[string]string{
				"core":   fmt.Sprintf("%d", i),
				"model":  cpu.ModelName,
				"vendor": cpu.VendorID,
			},
		})

		// one cpu same freq.
		if i == 0 {
			c.mu.Lock()
			c.stats.Frequency = freq
			c.mu.Unlock()
		}
	}

	metrics <- localMetrics
	errors <- nil
}

// func (c *cpuCollector) collectTemperature(now time.Time, metrics chan<- []Metric, errors chan<- error) {
// 	if !c.config.CollectTemperature {
// 		metrics <- nil
// 		errors <- nil
// 		return
// 	}
// 	// var localMetrics []Metric
// 	// TODO: Implement for MacOS, Windows and other OS
// }

func (c *cpuCollector) collectLoad(now time.Time, metrics chan<- []Metric, errors chan<- error) {
	if !c.config.CollectLoadAvg {
		metrics <- nil
		errors <- nil
		return
	}

	loadAvg, err := load.Avg()
	if err != nil {
		errors <- fmt.Errorf("load average: %v", err)
		return
	}

	numCPU := float64(runtime.NumCPU())
	perCPULoad := loadAvg.Load1 / numCPU

	// collect metrics for various time periods
	localMetrics := []Metric{
		{
			Name:      "system_load1",
			Value:     loadAvg.Load1,
			Timestamp: now,
			Labels:    map[string]string{"period": "1min"},
		},
		{
			Name:      "system_load5",
			Value:     loadAvg.Load5,
			Timestamp: now,
			Labels:    map[string]string{"period": "5min"},
		},
		{
			Name:      "system_load15",
			Value:     loadAvg.Load15,
			Timestamp: now,
			Labels:    map[string]string{"period": "15min"},
		},
		{
			Name:      "system_load_per_cpu",
			Value:     perCPULoad,
			Timestamp: now,
			Labels:    map[string]string{"period": "1min"},
		},
	}

	c.mu.Lock()
	c.stats.LoadAvg = perCPULoad
	c.loadHistory = append(c.loadHistory, perCPULoad)

	// replace if full
	if len(c.loadHistory) > c.config.HistorySize {
		c.loadHistory = c.loadHistory[1:]
	}
	c.mu.Unlock()

	if perCPULoad > c.config.LoadThreshold {
		c.addAlert(Alert{
			Source:    "system_load",
			Level:     AlertLevelWarning,
			Message:   fmt.Sprintf("High system load: %2f per CPU", perCPULoad),
			Timestamp: now,
		})
	}

	metrics <- localMetrics
	errors <- nil
}

func (c *cpuCollector) GetLoadTrend() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.loadHistory) < 2 {
		return "insufficient data"
	}

	first := c.loadHistory[0]
	last := c.loadHistory[len(c.loadHistory)-1]
	delta := last - first

	switch {
	case delta > 0.5:
		return "increasing rapidly"
	case delta > 0.1:
		return "increasing"
	case delta < -0.5:
		return "decreasing rapidly"
	case delta < -0.1:
		return "decreasing"
	default:
		return "stable"
	}
}

func (c *cpuCollector) GetCPUInfo() (map[string]interface{}, error) {
	info, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	// create k, v map with k[str], v[any value since interface]
	cpuInfo := make(map[string]interface{})

	if len(info) > 0 {
		cpuInfo["model"] = info[0].ModelName
		cpuInfo["cores"] = runtime.NumCPU()
		cpuInfo["vendor"] = info[0].VendorID
		cpuInfo["family"] = info[0].Family
		cpuInfo["stepping"] = info[0].Stepping
		cpuInfo["features"] = info[0].Flags
	}

	return cpuInfo, nil
}

func (c *cpuCollector) GetCPUStats() (*CPUStats, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats := c.stats
	return &stats, nil
}

func (c *cpuCollector) GetAlerts() []Alert {
	c.mu.Lock()
	defer c.mu.Unlock()

	alerts := c.alerts
	c.alerts = nil
	return alerts
}

func (c *cpuCollector) addAlert(alert Alert) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alerts = append(c.alerts, alert)
}
