package service

import (
	"context"
	"errors"
	"testing"

	"devops-monitoring-platform/backend/internal/model"
	"devops-monitoring-platform/backend/internal/repository"
)

type stubRepository struct {
	registerServerFn func(context.Context, model.ServerRegisterRequest) (*model.Server, error)
	listServersFn    func(context.Context) ([]model.Server, error)
	createMetricFn   func(context.Context, model.MetricCreateRequest) (*model.Metric, error)
	listMetricsFn    func(context.Context, int64) ([]model.Metric, error)
}

func (s *stubRepository) RegisterServer(ctx context.Context, req model.ServerRegisterRequest) (*model.Server, error) {
	return s.registerServerFn(ctx, req)
}

func (s *stubRepository) ListServers(ctx context.Context) ([]model.Server, error) {
	return s.listServersFn(ctx)
}

func (s *stubRepository) CreateMetric(ctx context.Context, req model.MetricCreateRequest) (*model.Metric, error) {
	return s.createMetricFn(ctx, req)
}

func (s *stubRepository) ListServerMetrics(ctx context.Context, serverID int64) ([]model.Metric, error) {
	return s.listMetricsFn(ctx, serverID)
}

func TestServerService_RegisterServer(t *testing.T) {
	repo := &stubRepository{
		registerServerFn: func(ctx context.Context, req model.ServerRegisterRequest) (*model.Server, error) {
			return &model.Server{Hostname: req.Hostname}, nil
		},
	}

	service := NewService(repo)
	server, err := service.RegisterServer(context.Background(), model.ServerRegisterRequest{Hostname: "web-01", IP: "10.0.0.1", OS: "Ubuntu"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if server.Hostname != "web-01" {
		t.Fatalf("expected hostname web-01, got %s", server.Hostname)
	}
}

func TestServerService_RegisterServerValidation(t *testing.T) {
	service := NewService(&stubRepository{})
	_, err := service.RegisterServer(context.Background(), model.ServerRegisterRequest{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestServerService_CreateMetric(t *testing.T) {
	repo := &stubRepository{
		createMetricFn: func(ctx context.Context, req model.MetricCreateRequest) (*model.Metric, error) {
			return &model.Metric{ServerID: req.ServerID}, nil
		},
	}
	service := NewService(repo)
	metric, err := service.CreateMetric(context.Background(), model.MetricCreateRequest{ServerID: 1, CPU: 10, Memory: 20, Disk: 30})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if metric.ServerID != 1 {
		t.Fatalf("expected server id 1, got %d", metric.ServerID)
	}
}

func TestServerService_ListServerMetrics(t *testing.T) {
	repo := &stubRepository{
		listMetricsFn: func(ctx context.Context, serverID int64) ([]model.Metric, error) {
			return []model.Metric{{ServerID: serverID}}, nil
		},
	}
	service := NewService(repo)
	metrics, err := service.ListServerMetrics(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metrics) != 1 {
		t.Fatalf("expected one metric, got %d", len(metrics))
	}
}

func TestServerService_ListServers(t *testing.T) {
	repo := &stubRepository{
		listServersFn: func(ctx context.Context) ([]model.Server, error) {
			return []model.Server{{Hostname: "web-01"}}, nil
		},
	}
	service := NewService(repo)
	servers, err := service.ListServers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(servers) != 1 {
		t.Fatalf("expected one server, got %d", len(servers))
	}
}

func TestServerService_PropagatesRepositoryError(t *testing.T) {
	repo := &stubRepository{
		listMetricsFn: func(ctx context.Context, serverID int64) ([]model.Metric, error) {
			return nil, repository.ErrServerNotFound
		},
	}
	service := NewService(repo)
	_, err := service.ListServerMetrics(context.Background(), 1)
	if !errors.Is(err, repository.ErrServerNotFound) {
		t.Fatalf("expected ErrServerNotFound, got %v", err)
	}
}
