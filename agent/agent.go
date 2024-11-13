package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"ddgo/internal/collector"

	"github.com/google/uuid"
)

// metrics collection agent
type Agent struct {
	ID              string
	Hostname        string
	ServerURL       string
	CPUCollector    *collector.CPUCollector
	MemoryCollector *collector.MemoryCollector
	DiskCollector   *collector.DiskCollector
}

// system metrics collected by agent
type AgentMetrics struct {
	AgentID  string `json:"agent_id"`
	Hostname string `json:"hostname"`
	Metrics  struct {
		CPU struct {
			Cores []struct {
				Core  int     `json:"core"`
				Usage float64 `json:"usage"`
			} `json:"cores"`
			Load struct {
				OneMin     float64 `json:"1m"`
				FiveMin    float64 `json:"5m"`
				FifteenMin float64 `json:"15m"`
			} `json:"load"`
			Times map[string]struct {
				User   float64 `json:"user"`
				System float64 `json:"system"`
				Idle   float64 `json:"idle"`
				IOWait float64 `json:"iowait"`
				IRQ    float64 `json:"irq"`
			} `json:"times"`
			Info struct {
				ProcessCount  int `json:"process_count"`
				ThreadCount   int `json:"thread_count"`
				LogicalCores  int `json:"logical_cores"`
				PhysicalCores int `json:"physical_cores"`
			} `json:"info"`
		} `json:"cpu"`
		Memory struct {
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
		} `json:"memory"`
		Disk struct {
			Usage float64 `json:"usage"`
			Total uint64  `json:"total"`
			Free  uint64  `json:"free"`
			IO    struct {
				ReadCount  uint64 `json:"read_count"`
				WriteCount uint64 `json:"write_count"`
				ReadBytes  uint64 `json:"read_bytes"`
				WriteBytes uint64 `json:"write_bytes"`
			} `json:"io"`
		} `json:"disk"`
		Time string `json:"time"`
	} `json:"metrics"`
	Timestamp time.Time `json:"timestamp"`
}

// create a new agent instance for server
func NewAgent(serverURL string) (*Agent, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %v", err)
	}

	return &Agent{
		ID:              uuid.New().String(),
		Hostname:        hostname,
		ServerURL:       serverURL,
		CPUCollector:    collector.CreateCPUCollector(0),
		MemoryCollector: collector.CreateMemoryCollector(),
		DiskCollector:   collector.CreateDiskCollector(),
	}, nil
}

// collect all metrics and send to server
func (a *Agent) CollectAndSend() error {
	metrics := AgentMetrics{
		AgentID:   a.ID,
		Hostname:  a.Hostname,
		Timestamp: time.Now(),
	}

	cpuMetrics, err := a.CPUCollector.Collect()
	if err != nil {
		return fmt.Errorf("CPU collection error: %v", err)
	}

	// parse cpu metrics
	for _, metric := range cpuMetrics {
		switch metric.Name {
		case "cpu_usage":
			if core, ok := metric.Labels["cpu"]; ok {
				var coreNum int
				fmt.Sscanf(core, "cpu%d", &coreNum)
				metrics.Metrics.CPU.Cores = append(metrics.Metrics.CPU.Cores, struct {
					Core  int     `json:"core"`
					Usage float64 `json:"usage"`
				}{
					Core:  coreNum,
					Usage: metric.Value,
				})
			}
		case "cpu_load_average_1m":
			metrics.Metrics.CPU.Load.OneMin = metric.Value
		case "cpu_load_average_5m":
			metrics.Metrics.CPU.Load.FiveMin = metric.Value
		case "cpu_load_average_15m":
			metrics.Metrics.CPU.Load.FifteenMin = metric.Value
		case "system_processes_total":
			metrics.Metrics.CPU.Info.ProcessCount = int(metric.Value)
		case "system_threads_total":
			metrics.Metrics.CPU.Info.ThreadCount = int(metric.Value)
		case "cpu_cores_logical":
			metrics.Metrics.CPU.Info.LogicalCores = int(metric.Value)
		case "cpu_cores_physical":
			metrics.Metrics.CPU.Info.PhysicalCores = int(metric.Value)
		}

		// cpu times
		if metric.Labels["cpu"] != "" {
			cpu := metric.Labels["cpu"]
			if metrics.Metrics.CPU.Times == nil {
				metrics.Metrics.CPU.Times = make(map[string]struct {
					User   float64 `json:"user"`
					System float64 `json:"system"`
					Idle   float64 `json:"idle"`
					IOWait float64 `json:"iowait"`
					IRQ    float64 `json:"irq"`
				})
			}
			times := metrics.Metrics.CPU.Times[cpu]
			switch metric.Name {
			case "cpu_time_user":
				times.User = metric.Value
			case "cpu_time_system":
				times.System = metric.Value
			case "cpu_time_idle":
				times.Idle = metric.Value
			case "cpu_time_iowait":
				times.IOWait = metric.Value
			case "cpu_time_irq":
				times.IRQ = metric.Value
			}
			metrics.Metrics.CPU.Times[cpu] = times
		}
	}

	// memory metrics
	memMetrics, err := a.MemoryCollector.Collect()
	if err != nil {
		return fmt.Errorf("Memory collection error: %v", err)
	}

	// parse memory metrics
	for _, metric := range memMetrics {
		if metric.Labels["type"] == "virtual" {
			switch metric.Name {
			case "memory_usage":
				metrics.Metrics.Memory.Virtual.Usage = metric.Value
			case "memory_total":
				metrics.Metrics.Memory.Virtual.Total = uint64(metric.Value)
			case "memory_used":
				metrics.Metrics.Memory.Virtual.Used = uint64(metric.Value)
			case "memory_free":
				metrics.Metrics.Memory.Virtual.Free = uint64(metric.Value)
			case "memory_cached":
				metrics.Metrics.Memory.Virtual.Cached = uint64(metric.Value)
			case "memory_available":
				metrics.Metrics.Memory.Virtual.Available = uint64(metric.Value)
			}
		} else if metric.Labels["type"] == "swap" {
			switch metric.Name {
			case "memory_usage":
				metrics.Metrics.Memory.Swap.Usage = metric.Value
			case "memory_total":
				metrics.Metrics.Memory.Swap.Total = uint64(metric.Value)
			case "memory_used":
				metrics.Metrics.Memory.Swap.Used = uint64(metric.Value)
			case "memory_free":
				metrics.Metrics.Memory.Swap.Free = uint64(metric.Value)
			}
		}
	}

	// disk metrics
	diskMetrics, err := a.DiskCollector.Collect()
	if err != nil {
		return fmt.Errorf("Disk collection error: %v", err)
	}

	// format disk metrics
	for _, metric := range diskMetrics {
		switch metric.Name {
		case "disk_usage":
			metrics.Metrics.Disk.Usage = metric.Value
		case "disk_total":
			metrics.Metrics.Disk.Total = uint64(metric.Value)
		case "disk_free":
			metrics.Metrics.Disk.Free = uint64(metric.Value)
		case "disk_reads_total":
			metrics.Metrics.Disk.IO.ReadCount = uint64(metric.Value)
		case "disk_writes_total":
			metrics.Metrics.Disk.IO.WriteCount = uint64(metric.Value)
		case "disk_read_bytes_total":
			metrics.Metrics.Disk.IO.ReadBytes = uint64(metric.Value)
		case "disk_write_bytes_total":
			metrics.Metrics.Disk.IO.WriteBytes = uint64(metric.Value)
		}
	}

	metrics.Metrics.Time = time.Now().Format(time.RFC3339)

	// marshal and send to server
	data, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/api/metrics/collect", a.ServerURL),
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return fmt.Errorf("failed to send metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}

	return nil
}

// start the agent and send metrics to server
func (a *Agent) Start() error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	log.Printf("Agent started. ID: %s, Hostname: %s", a.ID, a.Hostname)
	log.Printf("Sending metrics to: %s", a.ServerURL)

	for {
		select {
		case <-ticker.C:
			if err := a.CollectAndSend(); err != nil {
				log.Printf("Error collecting/sending metrics: %v", err)
			}
		}
	}
}
