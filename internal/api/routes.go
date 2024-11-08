package api

func (s *Server) routes() {
    s.mux.HandleFunc("/api/metrics", s.handleMetrics)
    s.mux.HandleFunc("/api/alerts", s.handleAlerts)
    s.mux.HandleFunc("/api/alerts/history", s.handleAlertHistory)
    s.mux.Handle("/", http.FileServer(http.Dir("web/dashboard")))
}