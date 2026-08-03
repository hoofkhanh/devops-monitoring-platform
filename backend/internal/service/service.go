package service

import (
	"context"

	"devops-monitoring-platform/backend/internal/model"
	"devops-monitoring-platform/backend/internal/repository"
)

type Service interface {
	RegisterServer(ctx context.Context, req model.ServerRegisterRequest) (*model.Server, error)
	ListServers(ctx context.Context) ([]model.Server, error)
	CreateMetric(ctx context.Context, req model.MetricCreateRequest) (*model.Metric, error)
	ListServerMetrics(ctx context.Context, serverID int64) ([]model.Metric, error)
}

type ServerService struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *ServerService {
	return &ServerService{repo: repo}
}

func (s *ServerService) RegisterServer(ctx context.Context, req model.ServerRegisterRequest) (*model.Server, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return s.repo.RegisterServer(ctx, req)
}

func (s *ServerService) ListServers(ctx context.Context) ([]model.Server, error) {
	return s.repo.ListServers(ctx)
}

func (s *ServerService) CreateMetric(ctx context.Context, req model.MetricCreateRequest) (*model.Metric, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return s.repo.CreateMetric(ctx, req)
}

func (s *ServerService) ListServerMetrics(ctx context.Context, serverID int64) ([]model.Metric, error) {
	return s.repo.ListServerMetrics(ctx, serverID)
}
