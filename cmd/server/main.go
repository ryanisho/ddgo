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

type DetailedMetrics struct {
	CPU struct {
		Cores []CoreMetric       `json:"cores"`
		Load  LoadMetrics        `json:"load"`
		Times map[string]CPUTime `json:"times"`
		Info  SystemInfo         `json:"info"`
	} `json:"cpu"`
	Memory MemoryMetrics `json:"memory"`
	Disk   DiskMetrics   `json:"disk"`
	Time   string        `json:"time"`
}

type CoreMetric struct {
	Core  int     `json:"core"`
	Usage float64 `json:"usage"`
}

type LoadMetrics struct {
	OneMin     float64 `json:"1m"`
	FiveMin    float64 `json:"5m"`
	FifteenMin float64 `json:"15m"`
}

type CPUTime struct {
	User   float64 `json:"user"`
	System float64 `json:"system"`
	Idle   float64 `json:"idle"`
	IOWait float64 `json:"iowait"`
	IRQ    float64 `json:"irq"`
}

type SystemInfo struct {
	ProcessCount  int `json:"process_count"`
	ThreadCount   int `json:"thread_count"`
	LogicalCores  int `json:"logical_cores"`
	PhysicalCores int `json:"physical_cores"`
}

type MemoryMetrics struct {
	Virtual struct {
		Total     uint64  `json:"total"`
		Used      uint64  `json:"used"`
		Free      uint64  `json:"free"`
		Usage     float64 `json:"usage"`
		Cached    uint64  `json:"cached"`
		Available uint64  `json:"available"`
	} `json:"virtual"`
	Swap struct {
		Total uint64  `json:"total"`
		Used  uint64  `json:"used"`
		Free  uint64  `json:"free"`
		Usage float64 `json:"usage"`
	} `json:"swap"`
}

type DiskMetrics struct {
	Usage float64 `json:"usage"`
	Total uint64  `json:"total"`
	Free  uint64  `json:"free"`
	IO    struct {
		ReadCount  uint64 `json:"read_count"`
		WriteCount uint64 `json:"write_count"`
		ReadBytes  uint64 `json:"read_bytes"`
		WriteBytes uint64 `json:"write_bytes"`
	} `json:"io"`
}

func collectMetrics() DetailedMetrics {
	var metrics DetailedMetrics
	metrics.Time = time.Now().Format(time.RFC3339)
	metrics.CPU.Times = make(map[string]CPUTime)

	var wg sync.WaitGroup
	wg.Add(3)

	// Collect CPU metrics
	go func() {
		defer wg.Done()
		cpuCollector := collector.CreateCPUCollector(0)
		if cpuMetrics, err := cpuCollector.Collect(); err == nil {
			for _, metric := range cpuMetrics {
				switch metric.Name {
				case "cpu_usage":
					if core, ok := metric.Labels["cpu"]; ok {
						var coreNum int
						if _, err := fmt.Sscanf(core, "cpu%d", &coreNum); err == nil {
							metrics.CPU.Cores = append(metrics.CPU.Cores, CoreMetric{
								Core:  coreNum,
								Usage: metric.Value,
							})
						}
					}
				case "cpu_load_average_1m":
					metrics.CPU.Load.OneMin = metric.Value
				case "cpu_load_average_5m":
					metrics.CPU.Load.FiveMin = metric.Value
				case "cpu_load_average_15m":
					metrics.CPU.Load.FifteenMin = metric.Value
				case "system_processes_total":
					metrics.CPU.Info.ProcessCount = int(metric.Value)
				case "system_threads_total":
					metrics.CPU.Info.ThreadCount = int(metric.Value)
				case "cpu_cores_logical":
					metrics.CPU.Info.LogicalCores = int(metric.Value)
				case "cpu_cores_physical":
					metrics.CPU.Info.PhysicalCores = int(metric.Value)
				}

				if metric.Labels["cpu"] != "" {
					cpu := metric.Labels["cpu"]
					switch metric.Name {
					case "cpu_time_user":
						if _, ok := metrics.CPU.Times[cpu]; !ok {
							metrics.CPU.Times[cpu] = CPUTime{}
						}
						temp := metrics.CPU.Times[cpu]
						temp.User = metric.Value
						metrics.CPU.Times[cpu] = temp
					case "cpu_time_system":
						if _, ok := metrics.CPU.Times[cpu]; !ok {
							metrics.CPU.Times[cpu] = CPUTime{}
						}
						temp := metrics.CPU.Times[cpu]
						temp.System = metric.Value
						metrics.CPU.Times[cpu] = temp
					case "cpu_time_idle":
						if _, ok := metrics.CPU.Times[cpu]; !ok {
							metrics.CPU.Times[cpu] = CPUTime{}
						}
						temp := metrics.CPU.Times[cpu]
						temp.Idle = metric.Value
						metrics.CPU.Times[cpu] = temp
					case "cpu_time_iowait":
						if _, ok := metrics.CPU.Times[cpu]; !ok {
							metrics.CPU.Times[cpu] = CPUTime{}
						}
						temp := metrics.CPU.Times[cpu]
						temp.IOWait = metric.Value
						metrics.CPU.Times[cpu] = temp
					case "cpu_time_irq":
						if _, ok := metrics.CPU.Times[cpu]; !ok {
							metrics.CPU.Times[cpu] = CPUTime{}
						}
						temp := metrics.CPU.Times[cpu]
						temp.IRQ = metric.Value
						metrics.CPU.Times[cpu] = temp
					}
				}
			}
		}
	}()

	// Collect Memory metrics
	go func() {
		defer wg.Done()
		memCollector := collector.CreateMemoryCollector()
		if memMetrics, err := memCollector.Collect(); err == nil {
			for _, metric := range memMetrics {
				if metric.Labels["type"] == "virtual" {
					switch metric.Name {
					case "memory_usage":
						metrics.Memory.Virtual.Usage = metric.Value
					case "memory_total":
						metrics.Memory.Virtual.Total = uint64(metric.Value)
					case "memory_used":
						metrics.Memory.Virtual.Used = uint64(metric.Value)
					case "memory_free":
						metrics.Memory.Virtual.Free = uint64(metric.Value)
					case "memory_cached":
						metrics.Memory.Virtual.Cached = uint64(metric.Value)
					case "memory_available":
						metrics.Memory.Virtual.Available = uint64(metric.Value)
					}
				} else if metric.Labels["type"] == "swap" {
					switch metric.Name {
					case "memory_usage":
						metrics.Memory.Swap.Usage = metric.Value
					case "memory_total":
						metrics.Memory.Swap.Total = uint64(metric.Value)
					case "memory_used":
						metrics.Memory.Swap.Used = uint64(metric.Value)
					case "memory_free":
						metrics.Memory.Swap.Free = uint64(metric.Value)
					}
				}
			}
		}
	}()

	// Collect Disk metrics
	go func() {
		defer wg.Done()
		diskCollector := collector.CreateDiskCollector()
		if diskMetrics, err := diskCollector.Collect(); err == nil {
			for _, metric := range diskMetrics {
				switch metric.Name {
				case "disk_usage":
					metrics.Disk.Usage = metric.Value
				case "disk_total":
					metrics.Disk.Total = uint64(metric.Value)
				case "disk_free":
					metrics.Disk.Free = uint64(metric.Value)
				case "disk_reads_total":
					metrics.Disk.IO.ReadCount = uint64(metric.Value)
				case "disk_writes_total":
					metrics.Disk.IO.WriteCount = uint64(metric.Value)
				case "disk_read_bytes_total":
					metrics.Disk.IO.ReadBytes = uint64(metric.Value)
				case "disk_write_bytes_total":
					metrics.Disk.IO.WriteBytes = uint64(metric.Value)
				}
			}
		}
	}()

	wg.Wait()
	return metrics
}

func main() {
	// Enable CORS
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	// Create router
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// API endpoint
	mux.HandleFunc("/api/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		metrics := collectMetrics()
		json.NewEncoder(w).Encode(metrics)
	})

	// Serve the frontend
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/index.html")
	})

	// Create server with CORS middleware
	handler := corsMiddleware(mux)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
