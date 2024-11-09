package storage

import (
	"sync"
	"time"
	"ddgo/internal/collector"
)


type MetricsStore struct {
	metrics []collector.Metric
	mu sync.RWMutex
}

func NewMetricsStore() *MetricsStore {
	return &MetricsStore{
		metrics: make([]collector.Metric, 0),
	}
}

func (s *MetricsStore) Store(metrics []collector.Metric) {
	s.mu.Lock() 
	defer s.mu.Unlock()

	s.metrics = append(s.metrics, metrics...)

	// last hour of metrics
	cutoff := time.Now().Add(-1 * time.Hour)
	filtered := make([]collector.Metric, 0)

	for _, m := range s.metrics {
		if m.Timestamp.After(cutoff) {
			filtered = append(filtered, m)
		}
	}
	s.metrics = filtered
}

func (s *MetricsStore) GetMetrics() []collector.Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.metrics
}