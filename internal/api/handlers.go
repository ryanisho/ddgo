package api

import (
	"encoding/json"
	"html/template"
	"net/http"
	"ddgo/internal/storage"
)

type Handler struct {
	store *storage.MetricsStore
}

func NewHandler(store *storage.MetricsStore) *Handler {
	return &Handler {store: store}
}

func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := h.store.GetMetrics()
	json.NewEncoder(w).Encode(metrics)
}

func (h *Handler) ServeUI(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("web/templates/index.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}