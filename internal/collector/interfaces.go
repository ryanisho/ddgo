package collector

import "time"

type CPUCollector interface {
	Collector
	GetCPUInfo() (map[string]interface{}, error)
	GetLoadTrend() string
	GetCPUStats() (*CPUStats, error)
}

type MemoryCollector interface {
	Collector
	GetMemoryStats() (*MemoryStats, error)
	GetSwapStats() (*MemoryStats, error)
}

type DiskCollector interface {
	Collector
	GetDiskStats(path string) (*DiskStats, error)
	GetIOStats(device string) (*DiskStats, error)
}

type MetricsStore interface {
	Store(metrics []Metric) error
	Query(name string, labels map[string]string, duration time.Duration) ([]Metric, error)
	Purge(age time.Duration) error
}

type AlertManager interface {
	AddAlert(alert Alert)
	GetAlerts() []Alert
	ClearAlerts()
	Subscribe(ch chan<- Alert)
	Unsubscribe(ch chan<- Alert)
}
