package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"devops-monitoring-platform/backend/internal/model"
)

var ErrServerNotFound = errors.New("server not found")

type Repository interface {
	RegisterServer(ctx context.Context, req model.ServerRegisterRequest) (*model.Server, error)
	ListServers(ctx context.Context) ([]model.Server, error)
	CreateMetric(ctx context.Context, req model.MetricCreateRequest) (*model.Metric, error)
	ListServerMetrics(ctx context.Context, serverID int64) ([]model.Metric, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) RegisterServer(ctx context.Context, req model.ServerRegisterRequest) (*model.Server, error) {
	query := `
		INSERT INTO servers (hostname, ip, os, status, last_seen)
		VALUES ($1, $2, $3, 'online', NOW())
		ON CONFLICT (hostname) DO UPDATE
		SET ip = EXCLUDED.ip,
		    os = EXCLUDED.os,
		    status = 'online',
		    last_seen = NOW()
		RETURNING id, hostname, ip, os, status, last_seen, created_at
	`

	var server model.Server
	var lastSeen sql.NullTime
	var createdAt time.Time

	err := r.db.QueryRowContext(ctx, query, req.Hostname, req.IP, req.OS).Scan(
		&server.ID,
		&server.Hostname,
		&server.IP,
		&server.OS,
		&server.Status,
		&lastSeen,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	if lastSeen.Valid {
		server.LastSeen = &lastSeen.Time
	}
	server.CreatedAt = createdAt

	return &server, nil
}

func (r *PostgresRepository) ListServers(ctx context.Context) ([]model.Server, error) {
	query := `
		SELECT id, hostname, ip, os, status, last_seen, created_at
		FROM servers
		ORDER BY last_seen DESC NULLS LAST, created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	servers := make([]model.Server, 0)
	for rows.Next() {
		var server model.Server
		var lastSeen sql.NullTime
		var createdAt time.Time

		err := rows.Scan(&server.ID, &server.Hostname, &server.IP, &server.OS, &server.Status, &lastSeen, &createdAt)
		if err != nil {
			return nil, err
		}
		if lastSeen.Valid {
			server.LastSeen = &lastSeen.Time
		}
		server.CreatedAt = createdAt
		servers = append(servers, server)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return servers, nil
}

func (r *PostgresRepository) CreateMetric(ctx context.Context, req model.MetricCreateRequest) (*model.Metric, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM servers WHERE id = $1)`, req.ServerID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrServerNotFound
	}

	metricQuery := `
		INSERT INTO metrics (server_id, cpu, memory, disk, timestamp)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, server_id, cpu, memory, disk, timestamp
	`

	var metric model.Metric
	var cpuRaw, memoryRaw, diskRaw []byte
	var timestamp time.Time

	err = tx.QueryRowContext(ctx, metricQuery, req.ServerID, req.CPU, req.Memory, req.Disk).Scan(
		&metric.ID,
		&metric.ServerID,
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

	_, err = tx.ExecContext(ctx, `UPDATE servers SET last_seen = NOW(), status = 'online' WHERE id = $1`, req.ServerID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &metric, nil
}

func (r *PostgresRepository) ListServerMetrics(ctx context.Context, serverID int64) ([]model.Metric, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM servers WHERE id = $1)`, serverID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrServerNotFound
	}

	query := `
		SELECT id, server_id, cpu, memory, disk, timestamp
		FROM metrics
		WHERE server_id = $1
		ORDER BY timestamp DESC, id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metrics := make([]model.Metric, 0)
	for rows.Next() {
		var metric model.Metric
		var cpuRaw, memoryRaw, diskRaw []byte
		var timestamp time.Time

		err := rows.Scan(&metric.ID, &metric.ServerID, &cpuRaw, &memoryRaw, &diskRaw, &timestamp)
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
