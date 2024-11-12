package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"ddgo/internal/collector"
)

type CoreMetric struct {
	Core  int     `json:"core"`
	Usage float64 `json:"usage"`
}

type Metrics struct {
	CPUCores []CoreMetric `json:"cpu_cores"`
	CPUTotal float64      `json:"cpu_total"`
	Memory   float64      `json:"memory"`
	Disk     float64      `json:"disk"`
	Time     string       `json:"time"`
}

func getMetrics() Metrics {
	cpuCollector := collector.CreateCPUCollector(0)
	memCollector := collector.CreateMemoryCollector()
	diskCollector := collector.CreateDiskCollector()

	cpuMetrics, err := cpuCollector.Collect()

	// log.Printf("===========: %v", cpuMetrics)

	if err != nil {
		log.Printf("Error collecting CPU metrics: %v", err)
	}

	memMetrics, err := memCollector.Collect()

	// log.Printf("===========: %v", memMetrics)
	if err != nil {
		log.Printf("Error collecting memory metrics: %v", err)
	}

	diskMetrics, err := diskCollector.Collect()
	log.Printf("===========: %v", diskMetrics)

	if err != nil {
		log.Printf("Error collecting disk metrics: %v", err)
	}

	var cpuCores []CoreMetric
	var cpuTotal float64

	for _, metric := range cpuMetrics {
		if metric.Name == "cpu_usage" {
			if core, ok := metric.Labels["cpu"]; ok {
				if core != "total" {
					coreNum := 0
					fmt.Sscanf(core, "cpu%d", &coreNum)
					cpuCores = append(cpuCores, CoreMetric{
						Core:  coreNum,
						Usage: metric.Value,
					})
				} else {
					cpuTotal = metric.Value
				}
			}
		}
	}

	var memoryUsage float64
	for _, metric := range memMetrics {
		if metric.Name == "memory_usage" && metric.Labels["type"] == "virtual" {
			memoryUsage = metric.Value
		}
	}

	var diskUsage float64
	for _, metric := range diskMetrics {
		if metric.Name == "disk_usage" {
			diskUsage = metric.Value
			break
		}
	}

	metrics := Metrics{
		CPUCores: cpuCores,
		CPUTotal: cpuTotal,
		Memory:   memoryUsage,
		Disk:     diskUsage,
		Time:     time.Now().Format(time.RFC3339),
	}

	return metrics
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// API endpoint
	http.HandleFunc("/api/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics := getMetrics()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/index.html")
	})

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
