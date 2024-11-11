package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
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
	// Get per-core CPU usage
	perCPU, _ := cpu.Percent(time.Second, true)
	totalCPU, _ := cpu.Percent(time.Second, false)

	// Convert to CoreMetric slice
	cores := make([]CoreMetric, len(perCPU))
	for i, usage := range perCPU {
		cores[i] = CoreMetric{
			Core:  i,
			Usage: usage,
		}
	}

	// Get memory usage
	memory, _ := mem.VirtualMemory()

	// Get disk usage
	disk, _ := disk.Usage("/")

	return Metrics{
		CPUCores: cores,
		CPUTotal: totalCPU[0],
		Memory:   memory.UsedPercent,
		Disk:     disk.UsedPercent,
		Time:     time.Now().Format(time.RFC3339),
	}
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

	// Serve index page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/index.html")
	})

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
