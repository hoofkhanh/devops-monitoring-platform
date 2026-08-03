package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"devops-monitoring-platform/backend/internal/model"
	"devops-monitoring-platform/backend/internal/repository"
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
	r.Post("/api/servers/register", h.RegisterServer)
	r.Post("/api/metrics", h.CreateMetric)
	r.Get("/api/servers", h.ListServers)
	r.Get("/api/servers/{id}/metrics", h.ListServerMetrics)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) RegisterServer(w http.ResponseWriter, r *http.Request) {
	var req model.ServerRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	server, err := h.service.RegisterServer(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, server)
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

func (h *Handler) ListServers(w http.ResponseWriter, r *http.Request) {
	servers, err := h.service.ListServers(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, servers)
}

func (h *Handler) ListServerMetrics(w http.ResponseWriter, r *http.Request) {
	serverIDParam := chi.URLParam(r, "id")
	serverID, err := strconv.ParseInt(serverIDParam, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid server id")
		return
	}

	metrics, err := h.service.ListServerMetrics(r.Context(), serverID)
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
	if errors.Is(err, repository.ErrServerNotFound) {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	var validationErr interface{ Error() string }
	if errors.As(err, &validationErr) {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeError(w, http.StatusInternalServerError, "internal server error")
}
