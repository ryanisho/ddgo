package collector

const (
	// Metric Names
	MetricCPUUsage    = "cpu_usage"
	MetricCPUTemp     = "cpu_temperature"
	MetricCPUFreq     = "cpu_frequency"
	MetricMemoryUsage = "memory_usage"
	MetricDiskUsage   = "disk_usage"
	MetricDiskIO      = "disk_io"

	// Alert Levels
	AlertLevelInfo     = "info"
	AlertLevelWarning  = "warning"
	AlertLevelCritical = "critical"

	// Default Configuration Values
	DefaultHistorySize    = 60
	DefaultUsageThreshold = 90.0
	DefaultTempThreshold  = 80.0
	DefaultLoadThreshold  = 1.5

	// Time Intervals
	DefaultCollectionInterval = 30 // seconds
	DefaultRetentionPeriod    = 24 // hours
)

var SkipFilesystems = map[string]bool{
	"devfs":    true,
	"tmpfs":    true,
	"devtmpfs": true,
	"squashfs": true,
	"iso9660":  true,
	// ... other filesystem types
}
