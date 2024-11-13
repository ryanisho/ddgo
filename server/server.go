// server/server.go
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// MetricsServer handles collecting metrics from agents
type MetricsServer struct {
	// Store metrics from all agents
	agents map[string]AgentMetrics
	mu     sync.RWMutex
	// Optional: add configuration fields here
}

// AgentMetrics matches the agent's metrics structure
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

func NewMetricsServer() *MetricsServer {
	return &MetricsServer{
		agents: make(map[string]AgentMetrics),
	}
}

// HandleMetricsCollection handles incoming metrics from agents
func (s *MetricsServer) HandleMetricsCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var metrics AgentMetrics
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		http.Error(w, fmt.Sprintf("Invalid metrics data: %v", err), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.agents[metrics.AgentID] = metrics
	s.mu.Unlock()

	log.Printf("Received metrics from agent %s (%s)", metrics.AgentID, metrics.Hostname)
}

// HandleGetMetrics returns metrics for all agents
func (s *MetricsServer) HandleGetMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.RLock()
	response := make(map[string]AgentMetrics)
	for id, metrics := range s.agents {
		response[id] = metrics
	}
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CleanupOldMetrics removes metrics from inactive agents
func (s *MetricsServer) CleanupOldMetrics() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		threshold := time.Now().Add(-5 * time.Minute)

		s.mu.Lock()
		for id, metrics := range s.agents {
			if metrics.Timestamp.Before(threshold) {
				delete(s.agents, id)
				log.Printf("Removed inactive agent: %s (%s)", id, metrics.Hostname)
			}
		}
		s.mu.Unlock()
	}
}
