package db

import (
	"testing"
)

func TestConnectMissingEnv(t *testing.T) {
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "")

	_, err := Connect()
	if err == nil {
		t.Fatal("expected error when environment variables are missing")
	}
}
