package model

import (
	"testing"
)

func TestServerRegisterRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		request ServerRegisterRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: ServerRegisterRequest{
				Hostname: "web-01",
				IP:       "192.168.1.10",
				OS:       "Ubuntu 22.04",
			},
			wantErr: false,
		},
		{
			name: "missing hostname",
			request: ServerRegisterRequest{
				IP: "192.168.1.10",
				OS: "Ubuntu 22.04",
			},
			wantErr: true,
		},
		{
			name: "invalid ip",
			request: ServerRegisterRequest{
				Hostname: "web-01",
				IP:       "not-an-ip",
				OS:       "Ubuntu 22.04",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr && err == nil {
				t.Fatalf("expected validation error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestMetricCreateRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		request MetricCreateRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: MetricCreateRequest{
				ServerID: 1,
				CPU:      45.5,
				Memory:   70.1,
				Disk:     80.2,
			},
			wantErr: false,
		},
		{
			name: "negative metric",
			request: MetricCreateRequest{
				ServerID: 1,
				CPU:      -1,
			},
			wantErr: true,
		},
		{
			name: "invalid server id",
			request: MetricCreateRequest{
				ServerID: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr && err == nil {
				t.Fatalf("expected validation error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}
