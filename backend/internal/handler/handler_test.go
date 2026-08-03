package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"devops-monitoring-platform/backend/internal/model"
	"github.com/go-chi/chi/v5"
)

type stubService struct {
	registerServerFn func(context.Context, model.ServerRegisterRequest) (*model.Server, error)
	listServersFn    func(context.Context) ([]model.Server, error)
	createMetricFn   func(context.Context, model.MetricCreateRequest) (*model.Metric, error)
	listMetricsFn    func(context.Context, int64) ([]model.Metric, error)
}

func (s *stubService) RegisterServer(ctx context.Context, req model.ServerRegisterRequest) (*model.Server, error) {
	return s.registerServerFn(ctx, req)
}

func (s *stubService) ListServers(ctx context.Context) ([]model.Server, error) {
	return s.listServersFn(ctx)
}

func (s *stubService) CreateMetric(ctx context.Context, req model.MetricCreateRequest) (*model.Metric, error) {
	return s.createMetricFn(ctx, req)
}

func (s *stubService) ListServerMetrics(ctx context.Context, serverID int64) ([]model.Metric, error) {
	return s.listMetricsFn(ctx, serverID)
}

func TestHandlerHealth(t *testing.T) {
	handler := NewHandler(&stubService{})
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandlerRegisterServer(t *testing.T) {
	service := &stubService{
		registerServerFn: func(ctx context.Context, req model.ServerRegisterRequest) (*model.Server, error) {
			return &model.Server{Hostname: req.Hostname}, nil
		},
	}
	handler := NewHandler(service)
	payload := bytes.NewBufferString(`{"hostname":"web-01","ip":"10.0.0.1","os":"Ubuntu"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/servers/register", payload)
	w := httptest.NewRecorder()

	handler.RegisterServer(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
	var response model.Server
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Hostname != "web-01" {
		t.Fatalf("expected hostname web-01, got %s", response.Hostname)
	}
}

func TestHandlerCreateMetric(t *testing.T) {
	service := &stubService{
		createMetricFn: func(ctx context.Context, req model.MetricCreateRequest) (*model.Metric, error) {
			return &model.Metric{ServerID: req.ServerID}, nil
		},
	}
	handler := NewHandler(service)
	payload := bytes.NewBufferString(`{"server_id":1,"cpu":1.2,"memory":2.3,"disk":3.4}`)
	r := httptest.NewRequest(http.MethodPost, "/api/metrics", payload)
	w := httptest.NewRecorder()

	handler.CreateMetric(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestHandlerListServers(t *testing.T) {
	service := &stubService{
		listServersFn: func(ctx context.Context) ([]model.Server, error) {
			return []model.Server{{Hostname: "web-01"}}, nil
		},
	}
	handler := NewHandler(service)
	r := httptest.NewRequest(http.MethodGet, "/api/servers", nil)
	w := httptest.NewRecorder()

	handler.ListServers(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandlerListServerMetrics(t *testing.T) {
	service := &stubService{
		listMetricsFn: func(ctx context.Context, serverID int64) ([]model.Metric, error) {
			return []model.Metric{{ServerID: serverID}}, nil
		},
	}
	handler := NewHandler(service)
	r := httptest.NewRequest(http.MethodGet, "/api/servers/1/metrics", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.ListServerMetrics(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
