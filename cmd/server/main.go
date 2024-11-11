package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"ddgo/internal/collector"
)

type MetricsCollector struct {
	collectors     map[string]collector.Collector
	config         collector.CollectorConfig
	store          collector.MetricsStore
	alerts         collector.AlertManager
	mutex          sync.RWMutex
	metricsHistory map[string][]collector.Metric
	lastCollection time.Time
}

type DetailedCPUStats struct {
	PerCore     []CoreStats      `json:"per_core"`
	Total       TotalStats       `json:"total"`
	Load        LoadStats        `json:"load"`
	Temperature []float64        `json:"temperature"`
	Frequency   []FrequencyStats `json:"frequency"`
	Trends      TrendStats       `json:"trends"`
	Times       []TimeStats      `json:"times"`
}

type CoreStats struct {
	CoreID      int     `json:"core_id"`
	Usage       float64 `json:"usage"`
	Frequency   float64 `json:"frequency"`
	Temperature float64 `json:"temperature,omitempty"`
}

type TotalStats struct {
	Usage       float64 `json:"usage"`
	Temperature float64 `json:"temperature"`
}

type LoadStats struct {
	OneMin     float64 `json:"1min"`
	FiveMin    float64 `json:"5min"`
	FifteenMin float64 `json:"15min"`
	PerCPU     float64 `json:"per_cpu"`
}

type FrequencyStats struct {
	CoreID  int     `json:"core_id"`
	Current float64 `json:"current"`
	Min     float64 `json:"min,omitempty"`
	Max     float64 `json:"max,omitempty"`
}

type TrendStats struct {
	LoadTrend  string    `json:"load_trend"`
	UsageTrend string    `json:"usage_trend"`
	LastUpdate time.Time `json:"last_update"`
}

type TimeStats struct {
	CoreID int     `json:"core_id"`
	User   float64 `json:"user"`
	System float64 `json:"system"`
	Idle   float64 `json:"idle"`
	IOWait float64 `json:"iowait,omitempty"`
}

func NewMetricsCollector() *MetricsCollector {
	config := collector.CollectorConfig{
		CPU: collector.CPUMetricConfig{
			CollectTemperature: true,
			CollectFrequency:   true,
			CollectLoadAvg:     true,
			HistorySize:        60,
			UsageThreshold:     90.0,
			TempThreshold:      80.0,
			LoadThreshold:      1.5,
		},
	}

	return &MetricsCollector{
		collectors:     make(map[string]collector.Collector),
		config:         config,
		metricsHistory: make(map[string][]collector.Metric),
		lastCollection: time.Now(),
	}
}

func (mc *MetricsCollector) collectDetailedCPUStats() (*DetailedCPUStats, error) {
	cpuCollector, ok := mc.collectors["cpu"].(collector.CPUCollector)
	if !ok {
		return nil, fmt.Errorf("CPU collector not properly initialized")
	}

	// Get CPU info
	cpuInfo, err := cpuCollector.GetCPUInfo()
	if err != nil {
		return nil, fmt.Errorf("error getting CPU info: %v", err)
	}

	// Log CPU info
	log.Printf("CPU Info: Model=%s, Cores=%d, Vendor=%s",
		cpuInfo["model"], cpuInfo["cores"], cpuInfo["vendor"])

	// Get CPU stats
	cpuStats, err := cpuCollector.GetCPUStats()
	if err != nil {
		return nil, fmt.Errorf("error getting CPU stats: %v", err)
	}

	// Log CPU stats
	log.Printf("CPU Stats: Usage=%.2f%%, Temperature=%.2f°C, Frequency=%.2fMHz, LoadAvg=%.2f",
		cpuStats.Usage, cpuStats.Temperature, cpuStats.Frequency, cpuStats.LoadAvg)

	// Get all metrics
	metrics, err := cpuCollector.Collect()
	if err != nil {
		return nil, err
	}

	// Get load trend
	loadTrend := cpuCollector.GetLoadTrend()

	// Process all collected data into DetailedCPUStats
	stats := &DetailedCPUStats{
		PerCore: make([]CoreStats, 0),
		Total: TotalStats{
			Usage:       cpuStats.Usage,
			Temperature: cpuStats.Temperature,
		},
		Load:        LoadStats{},
		Temperature: make([]float64, 0),
		Frequency: []FrequencyStats{{
			CoreID:  0,
			Current: cpuStats.Frequency,
		}},
		Trends: TrendStats{
			LoadTrend:  loadTrend,
			UsageTrend: fmt.Sprintf("%.2f%%", cpuStats.Usage),
			LastUpdate: time.Now(),
		},
	}

	// Add CPU info to the stats
	if vendor, ok := cpuInfo["vendor"].(string); ok {
		log.Printf("CPU Vendor: %s", vendor)
	}
	if model, ok := cpuInfo["model"].(string); ok {
		log.Printf("CPU Model: %s", model)
	}
	if cores, ok := cpuInfo["cores"].(int); ok {
		log.Printf("CPU Cores: %d", cores)
	}
	if features, ok := cpuInfo["features"].([]string); ok {
		log.Printf("CPU Features: %v", features)
	}

	// Process each metric
	for _, metric := range metrics {
		switch metric.Name {
		case "cpu_usage":
			if metric.Labels["type"] == "core" {
				coreID := 0
				fmt.Sscanf(metric.Labels["core"], "%d", &coreID)
				stats.PerCore = append(stats.PerCore, CoreStats{
					CoreID: coreID,
					Usage:  metric.Value,
				})
			}
		case "system_load1":
			stats.Load.OneMin = metric.Value
		case "system_load5":
			stats.Load.FiveMin = metric.Value
		case "system_load15":
			stats.Load.FifteenMin = metric.Value
		case "system_load_per_cpu":
			stats.Load.PerCPU = metric.Value
		}
	}

	return stats, nil
}

func main() {
	mc := NewMetricsCollector()

	// Initialize CPU collector with all features enabled
	cpuConfig := collector.CPUMetricConfig{
		CollectTemperature: true,
		CollectFrequency:   true,
		CollectLoadAvg:     true,
		HistorySize:        60,
		UsageThreshold:     90.0,
		TempThreshold:      80.0,
		LoadThreshold:      1.5,
	}

	mc.collectors["cpu"] = collector.CreateParameterizedCPUCollector(cpuConfig)

	// Start collection routine with enhanced logging
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if cpuCollector, ok := mc.collectors["cpu"].(collector.CPUCollector); ok {
				// Get CPU info periodically
				cpuInfo, err := cpuCollector.GetCPUInfo()
				if err != nil {
					log.Printf("Error getting CPU info: %v", err)
				} else {
					log.Printf("CPU Info: %+v", cpuInfo)
				}

				// Get CPU stats
				cpuStats, err := cpuCollector.GetCPUStats()
				if err != nil {
					log.Printf("Error getting CPU stats: %v", err)
				} else {
					log.Printf("CPU Stats: Usage=%.2f%%, Temp=%.2f°C, Freq=%.2fMHz, Load=%.2f",
						cpuStats.Usage, cpuStats.Temperature, cpuStats.Frequency, cpuStats.LoadAvg)
				}

				// Get detailed stats
				stats, err := mc.collectDetailedCPUStats()
				if err != nil {
					log.Printf("Error collecting detailed CPU stats: %v", err)
					continue
				}

				// Get and process alerts
				alerts := cpuCollector.GetAlerts()
				for _, alert := range alerts {
					log.Printf("[%s] %s: %s", alert.Level, alert.Source, alert.Message)
				}

				// Log trends and significant changes
				log.Printf("CPU Load Trend: %s", stats.Trends.LoadTrend)
				log.Printf("Per-Core Stats:")
				for _, core := range stats.PerCore {
					log.Printf("  Core %d: Usage=%.2f%%", core.CoreID, core.Usage)
				}
			}
		}
	}()

	// API endpoints remain the same...
	http.HandleFunc("/api/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics, err := mc.collectDetailedCPUStats()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics)
	})

	http.HandleFunc("/api/cpu/stats", func(w http.ResponseWriter, r *http.Request) {
		stats, err := mc.collectDetailedCPUStats()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})

	// Rest of the main function remains the same...
	http.HandleFunc("/api/cpu/trend", func(w http.ResponseWriter, r *http.Request) {
		if cpuCollector, ok := mc.collectors["cpu"].(collector.CPUCollector); ok {
			trend := cpuCollector.GetLoadTrend()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"trend": trend})
		}
	})

	// Start collection routine
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if cpuCollector, ok := mc.collectors["cpu"].(collector.CPUCollector); ok {
				stats, err := mc.collectDetailedCPUStats()
				if err != nil {
					log.Printf("Error collecting CPU stats: %v", err)
					continue
				}

				// Get alerts
				alerts := cpuCollector.GetAlerts()
				for _, alert := range alerts {
					log.Printf("[%s] %s: %s", alert.Level, alert.Source, alert.Message)
				}

				// Log trends and significant changes
				log.Printf("CPU Load Trend: %s", stats.Trends.LoadTrend)

				if stats.Total.Usage > mc.config.CPU.UsageThreshold {
					log.Printf("High CPU usage detected: %.2f%%", stats.Total.Usage)
				}
			}
		}
	}()

	// Serve static files and index
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/index.html")
	})

	log.Printf("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
