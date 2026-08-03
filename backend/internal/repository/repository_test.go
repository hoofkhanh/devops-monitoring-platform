package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"devops-monitoring-platform/backend/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgresRepository_RegisterServer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	mock.ExpectQuery("INSERT INTO servers").
		WithArgs("web-01", "10.0.0.1", "Ubuntu").
		WillReturnRows(sqlmock.NewRows([]string{"id", "hostname", "ip", "os", "status", "last_seen", "created_at"}).AddRow(1, "web-01", "10.0.0.1", "Ubuntu", "online", time.Now(), time.Now()))

	server, err := repo.RegisterServer(context.Background(), model.ServerRegisterRequest{Hostname: "web-01", IP: "10.0.0.1", OS: "Ubuntu"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if server.Hostname != "web-01" {
		t.Fatalf("expected hostname web-01, got %s", server.Hostname)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPostgresRepository_ListServers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	mock.ExpectQuery("SELECT id, hostname, ip, os, status, last_seen, created_at").
		WillReturnRows(sqlmock.NewRows([]string{"id", "hostname", "ip", "os", "status", "last_seen", "created_at"}).AddRow(1, "web-01", "10.0.0.1", "Ubuntu", "online", time.Now(), time.Now()))

	servers, err := repo.ListServers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(servers) != 1 {
		t.Fatalf("expected one server, got %d", len(servers))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPostgresRepository_CreateMetric(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM servers WHERE id = \$1\)`).WithArgs(int64(1)).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectQuery("INSERT INTO metrics").WithArgs(int64(1), 45.5, 60.5, 70.5).WillReturnRows(sqlmock.NewRows([]string{"id", "server_id", "cpu", "memory", "disk", "timestamp"}).AddRow(1, 1, "45.50", "60.50", "70.50", time.Now()))
	mock.ExpectExec("UPDATE servers SET last_seen").WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	metric, err := repo.CreateMetric(context.Background(), model.MetricCreateRequest{ServerID: 1, CPU: 45.5, Memory: 60.5, Disk: 70.5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if metric.ServerID != 1 {
		t.Fatalf("expected server id 1, got %d", metric.ServerID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPostgresRepository_ListServerMetrics(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM servers WHERE id = \$1\)`).WithArgs(int64(1)).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectQuery("SELECT id, server_id, cpu, memory, disk, timestamp").WithArgs(int64(1)).WillReturnRows(sqlmock.NewRows([]string{"id", "server_id", "cpu", "memory", "disk", "timestamp"}).AddRow(1, 1, "45.50", "60.50", "70.50", time.Now()))

	metrics, err := repo.ListServerMetrics(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metrics) != 1 {
		t.Fatalf("expected one metric, got %d", len(metrics))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestParseNumeric(t *testing.T) {
	value, err := parseNumeric([]byte("12.34"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != 12.34 {
		t.Fatalf("expected 12.34, got %v", value)
	}
}

func TestErrServerNotFound(t *testing.T) {
	if errors.Is(ErrServerNotFound, sql.ErrNoRows) {
		t.Fatal("expected a distinct not found error")
	}
}
