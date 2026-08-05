package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"devops-monitoring-platform/backend/internal/model"
)

func TestPostgresRepositoryCreateMetric(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sql mock: %v", err)
	}
	defer db.Close()

	timestamp := time.Now()
	mock.ExpectQuery("INSERT INTO metrics").
		WithArgs(45.5, 60.5, 70.5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "cpu", "memory", "disk", "timestamp"}).AddRow(1, "45.50", "60.50", "70.50", timestamp))

	repo := NewPostgresRepository(db)
	metric, err := repo.CreateMetric(context.Background(), model.MetricCreateRequest{CPU: 45.5, Memory: 60.5, Disk: 70.5})
	if err != nil {
		t.Fatalf("CreateMetric() error = %v", err)
	}
	if metric.CPU != 45.5 || metric.Memory != 60.5 || metric.Disk != 70.5 {
		t.Fatalf("unexpected metric: %+v", metric)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresRepositoryListMetrics(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("new sql mock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, cpu, memory, disk, timestamp").
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"id", "cpu", "memory", "disk", "timestamp"}).AddRow(1, "45.50", "60.50", "70.50", time.Now()))

	repo := NewPostgresRepository(db)
	metrics, err := repo.ListMetrics(context.Background())
	if err != nil {
		t.Fatalf("ListMetrics() error = %v", err)
	}
	if len(metrics) != 1 || metrics[0].ID != 1 {
		t.Fatalf("unexpected metrics: %+v", metrics)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
