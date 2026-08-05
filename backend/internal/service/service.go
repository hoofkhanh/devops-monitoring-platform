package service

import (
	"context"

	"devops-monitoring-platform/backend/internal/model"
	"devops-monitoring-platform/backend/internal/repository"
)

type Service interface {
	CreateMetric(ctx context.Context, req model.MetricCreateRequest) (*model.Metric, error)
	ListMetrics(ctx context.Context) ([]model.Metric, error)
}

type MetricsService struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *MetricsService {
	return &MetricsService{repo: repo}
}

func (s *MetricsService) CreateMetric(ctx context.Context, req model.MetricCreateRequest) (*model.Metric, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return s.repo.CreateMetric(ctx, req)
}

func (s *MetricsService) ListMetrics(ctx context.Context) ([]model.Metric, error) {
	return s.repo.ListMetrics(ctx)
}
