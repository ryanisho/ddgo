package collector

import "time"

type Metric struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
	Labels    map[string]string `json:"labels"`
}

type CPUStats struct {
	Usage       float64 `json:"usage"`
	Temperature float64 `json:"temperature"`
	Frequency   float64 `json:"frequency"`
	LoadAvg     float64 `json:"load_avg"`
}

type CPUMetricConfig struct {
	CollectTemperature bool    `json:"collect_temperature"`
	CollectFrequency   bool    `json:"collect_frequency"`
	CollectLoadAvg     bool    `json:"collect_load_avg"`
	HistorySize        int     `json:"history_size"`
	UsageThreshold     float64 `json:"usage_threshold"`
	TempThreshold      float64 `json:"temp_threshold"`
	LoadThreshold      float64 `json:"load_threshold"`
}

type MemoryStats struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
	SwapTotal   uint64  `json:"swap_total"`
	SwapUsed    uint64  `json:"swap_used"`
	SwapPercent float64 `json:"swap_percent"`
}

type DiskStats struct {
	Path        string  `json:"path"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
	ReadRate    float64 `json:"read_rate"`
	WriteRate   float64 `json:"write_rate"`
	IOPSRead    float64 `json:"iops_read"`
	IOPSWrite   float64 `json:"iops_write"`
}

type CollectorConfig struct {
	CPU    CPUMetricConfig `json:"cpu"`
	Memory struct {
		CollectSwap bool `json:"collect_swap"`
		HistorySize int  `json:"history_size"`
	} `json:"memory"`
	Disk struct {
		IncludePaths []string `json:"include_paths"`
		ExcludePaths []string `json:"exclude_paths"`
		CollectIO    bool     `json:"collect_io"`
		HistorySize  int      `json:"history_size"`
	} `json:"disk"`
}

type Alert struct {
	Source    string    `json:"source"`
	Level     string    `json:"level"` // "info", "warning", "critical"
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type Collector interface {
	Collect() ([]Metric, error)
	GetAlerts() []Alert
}
