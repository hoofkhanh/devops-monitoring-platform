package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"devops-monitoring-platform/backend/internal/model"
	"devops-monitoring-platform/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.Health)
	r.Post("/api/metrics", h.CreateMetric)
	r.Get("/api/metrics", h.ListMetrics)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) CreateMetric(w http.ResponseWriter, r *http.Request) {
	var req model.MetricCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	metric, err := h.service.CreateMetric(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, metric)
}

func (h *Handler) ListMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.service.ListMetrics(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, metrics)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	var validationErr interface{ Error() string }
	if errors.As(err, &validationErr) {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeError(w, http.StatusInternalServerError, "internal server error")
}
