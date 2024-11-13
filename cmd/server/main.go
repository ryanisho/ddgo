// cmd/server/main.go
package main

import (
	"flag"
	"log"
	"net/http"

	"ddgo/server"
)

func main() {
	port := flag.String("port", "8080", "Port to run the server on")
	flag.Parse()

	// Create metrics server
	metricsServer := server.NewMetricsServer()

	// Start cleanup routine
	go metricsServer.CleanupOldMetrics()

	// Setup routes
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/metrics/collect", metricsServer.HandleMetricsCollection)
	mux.HandleFunc("/api/metrics", metricsServer.HandleGetMetrics)

	// Serve static files (your frontend)
	mux.Handle("/", http.FileServer(http.Dir("web")))

	// Add CORS middleware
	handler := enableCORS(mux)

	// Start server
	addr := ":" + *port
	log.Printf("Server starting on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}

func enableCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow your React app's origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
