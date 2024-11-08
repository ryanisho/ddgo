package api

import (
    "encoding/json"
    "net/http"
    "time"
)

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetMetrics(w, r)
	case http.MethodPost:
		s.handlePostMetric(w, r)
	default: 
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleGetMetrics(w http.ResponseWriter, r *http.Reqest) {
	name := r.URL.Query().Get("name")

	if name == "" {
		http.Error(w, "Metric name is required", http.StatusBadRequest)
		return
	}

	end := time.Now()
	start := end.Add(-1 * time.Hour)

	if startStr := r.URL.Query().Get("start"); startStr != "" {
		if t, err := time.Parse(time.RGC3339, startStr); err == nil {
			start = t
		}
	}

	if endStr := r.URL.Query().Get("end"); endStr != "" {
		if t, err := time.Parse(time.RGC3339, endStr); err == nil {
			end = t
		}
	}

	metrics, err := s.db.GetMetrics(name, start, end)
	if err != nil {
		http.Error (w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(metrics)
}

func (s *Server) handlePostMetric(w http.ResponseWriter, r *http.Request) {
    var metric db.Metric
    if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    metric.Timestamp = time.Now()
    if err := s.db.SaveMetric(&metric); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleAlerts(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        alerts, err := s.db.GetAlerts()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(alerts)

    case http.MethodPost:
        var alert db.Alert
        if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        if err := s.db.SaveAlert(&alert); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)

    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func (s *Server) handleAlertHistory(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    history, err := s.db.GetAlertHistory()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(history)
}