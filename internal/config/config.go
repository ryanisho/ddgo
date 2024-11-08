package config 

import (
	"os"
	"time"
)

type Config struct {
	DBPath string
	MetricRetention time.Duration
	CollectInterval time.Duration
	ServerPort string
}

func LoadConfig() *Config {
	return &Config {
		DBPath: getEnv("DB_PATH", "data/monitoring.db"),
		MetricsRetention: getDurationEnv("METRICS_RETENTION", 7 * 24 * time.Hour),
		CollectInterval: getDurationEnv("COLLECT_INTERVAL", 10 * time.Second),
		ServerPort: getEnv("SERVER_PORT", ":8080"),	

	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}