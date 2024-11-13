package main

import (
	"flag"
	"log"
	"net/http"

	"ddgo/server"
)

// CORS headers for response
func startCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func main() {
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	metricsServer := server.StartServer()
	go metricsServer.Clean() // flush old metrics every 5 minutes

	mux := http.NewServeMux() // routes

	// endpoints
	mux.HandleFunc("/api/metrics/collect", metricsServer.CollectAgents)
	mux.HandleFunc("/api/metrics", metricsServer.GetMetrics)

	handler := startCORS(mux)

	addr := ":" + *port // listen on all ports
	log.Printf("Server starting on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, handler)) // start server
}
