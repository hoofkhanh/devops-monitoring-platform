package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"devops-monitoring-platform/backend/internal/model"
	"github.com/go-chi/chi/v5"
)

type stubService struct {
	createMetricFn func(context.Context, model.MetricCreateRequest) (*model.Metric, error)
	listMetricsFn  func(context.Context) ([]model.Metric, error)
}

func (s *stubService) CreateMetric(ctx context.Context, req model.MetricCreateRequest) (*model.Metric, error) {
	return s.createMetricFn(ctx, req)
}

func (s *stubService) ListMetrics(ctx context.Context) ([]model.Metric, error) {
	return s.listMetricsFn(ctx)
}

func TestHandlerCreateMetric(t *testing.T) {
	h := NewHandler(&stubService{
		createMetricFn: func(_ context.Context, req model.MetricCreateRequest) (*model.Metric, error) {
			return &model.Metric{CPU: req.CPU}, nil
		},
	})
	w := httptest.NewRecorder()
	h.CreateMetric(w, httptest.NewRequest(http.MethodPost, "/api/metrics", strings.NewReader(`{"cpu":1.2,"memory":2.3,"disk":3.4}`)))
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestHandlerListMetrics(t *testing.T) {
	h := NewHandler(&stubService{
		listMetricsFn: func(context.Context) ([]model.Metric, error) { return []model.Metric{{CPU: 1.2}}, nil },
	})
	w := httptest.NewRecorder()
	h.ListMetrics(w, httptest.NewRequest(http.MethodGet, "/api/metrics", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"cpu":1.2`) {
		t.Fatalf("unexpected response: status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestHandlerRegistersOnlyMetricsRoutes(t *testing.T) {
	h := NewHandler(&stubService{
		listMetricsFn: func(context.Context) ([]model.Metric, error) { return []model.Metric{}, nil },
	})
	router := chi.NewRouter()
	h.RegisterRoutes(router)

	metricsResponse := httptest.NewRecorder()
	router.ServeHTTP(metricsResponse, httptest.NewRequest(http.MethodGet, "/api/metrics", nil))
	if metricsResponse.Code != http.StatusOK {
		t.Fatalf("expected metrics route status %d, got %d", http.StatusOK, metricsResponse.Code)
	}

	serverResponse := httptest.NewRecorder()
	router.ServeHTTP(serverResponse, httptest.NewRequest(http.MethodGet, "/api/servers", nil))
	if serverResponse.Code != http.StatusNotFound {
		t.Fatalf("expected removed servers route status %d, got %d", http.StatusNotFound, serverResponse.Code)
	}
}
