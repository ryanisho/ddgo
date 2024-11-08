package api 

import (
	"net/http"
	"ddgo/internal/collector"
	"ddgo/internal/db"
)

type Server struct {
	db *db.DB
	collector *collector.SystemCollector
	mux *http.ServeMux
}

func NewServer(database *db.DB, collector *collector.SystemCollector) *Server {
	server := &Server{
		db: database,
		collector: collector,
		mux: http.NewServeMux(),
	}
	
	server.routes()
	return server
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

