package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"devops-monitoring-platform/backend/internal/model"
)

type Repository interface {
	CreateMetric(ctx context.Context, req model.MetricCreateRequest) (*model.Metric, error)
	ListMetrics(ctx context.Context) ([]model.Metric, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateMetric(ctx context.Context, req model.MetricCreateRequest) (*model.Metric, error) {
	metricQuery := `
		INSERT INTO metrics (cpu, memory, disk, timestamp)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, cpu, memory, disk, timestamp
	`

	var metric model.Metric
	var cpuRaw, memoryRaw, diskRaw []byte
	var timestamp time.Time

	err := r.db.QueryRowContext(ctx, metricQuery, req.CPU, req.Memory, req.Disk).Scan(
		&metric.ID,
		&cpuRaw,
		&memoryRaw,
		&diskRaw,
		&timestamp,
	)
	if err != nil {
		return nil, err
	}

	metric.CPU, err = parseNumeric(cpuRaw)
	if err != nil {
		return nil, fmt.Errorf("parse cpu: %w", err)
	}
	metric.Memory, err = parseNumeric(memoryRaw)
	if err != nil {
		return nil, fmt.Errorf("parse memory: %w", err)
	}
	metric.Disk, err = parseNumeric(diskRaw)
	if err != nil {
		return nil, fmt.Errorf("parse disk: %w", err)
	}
	metric.Timestamp = timestamp

	return &metric, nil
}

func (r *PostgresRepository) ListMetrics(ctx context.Context) ([]model.Metric, error) {
	query := `
		SELECT id, cpu, memory, disk, timestamp
		FROM metrics
		ORDER BY timestamp DESC, id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metrics := make([]model.Metric, 0)
	for rows.Next() {
		var metric model.Metric
		var cpuRaw, memoryRaw, diskRaw []byte
		var timestamp time.Time

		err := rows.Scan(&metric.ID, &cpuRaw, &memoryRaw, &diskRaw, &timestamp)
		if err != nil {
			return nil, err
		}

		metric.CPU, err = parseNumeric(cpuRaw)
		if err != nil {
			return nil, fmt.Errorf("parse cpu: %w", err)
		}
		metric.Memory, err = parseNumeric(memoryRaw)
		if err != nil {
			return nil, fmt.Errorf("parse memory: %w", err)
		}
		metric.Disk, err = parseNumeric(diskRaw)
		if err != nil {
			return nil, fmt.Errorf("parse disk: %w", err)
		}
		metric.Timestamp = timestamp
		metrics = append(metrics, metric)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}

func parseNumeric(raw []byte) (float64, error) {
	if len(raw) == 0 {
		return 0, nil
	}
	return strconv.ParseFloat(string(raw), 64)
}
