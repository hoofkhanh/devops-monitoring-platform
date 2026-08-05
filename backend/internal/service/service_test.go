package service

import (
	"context"
	"errors"
	"testing"

	"devops-monitoring-platform/backend/internal/model"
)

type stubRepository struct {
	createMetricFn func(context.Context, model.MetricCreateRequest) (*model.Metric, error)
	listMetricsFn  func(context.Context) ([]model.Metric, error)
}

func (s *stubRepository) CreateMetric(ctx context.Context, req model.MetricCreateRequest) (*model.Metric, error) {
	return s.createMetricFn(ctx, req)
}

func (s *stubRepository) ListMetrics(ctx context.Context) ([]model.Metric, error) {
	return s.listMetricsFn(ctx)
}

func TestMetricsServiceCreateMetric(t *testing.T) {
	service := NewService(&stubRepository{
		createMetricFn: func(_ context.Context, req model.MetricCreateRequest) (*model.Metric, error) {
			return &model.Metric{CPU: req.CPU}, nil
		},
	})

	metric, err := service.CreateMetric(context.Background(), model.MetricCreateRequest{CPU: 10, Memory: 20, Disk: 30})
	if err != nil || metric.CPU != 10 {
		t.Fatalf("CreateMetric() = %+v, %v", metric, err)
	}
}

func TestMetricsServiceRejectsInvalidMetric(t *testing.T) {
	service := NewService(&stubRepository{})
	if _, err := service.CreateMetric(context.Background(), model.MetricCreateRequest{CPU: -1}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestMetricsServiceListMetrics(t *testing.T) {
	wantErr := errors.New("database unavailable")
	service := NewService(&stubRepository{
		listMetricsFn: func(context.Context) ([]model.Metric, error) { return nil, wantErr },
	})
	if _, err := service.ListMetrics(context.Background()); !errors.Is(err, wantErr) {
		t.Fatalf("ListMetrics() error = %v, want %v", err, wantErr)
	}
}
