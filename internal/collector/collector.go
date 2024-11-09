package collector

import (
	"time"
)

// metric obj.
type Metric struct {
	Name string
	Value float64 
	Timestamp time.Time
	Labels map[string] string
}

type Collector interface {
	Collect() ([]Metric, error)
}