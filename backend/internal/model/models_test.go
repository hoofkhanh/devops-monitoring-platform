package model

import "testing"

func TestMetricCreateRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		request MetricCreateRequest
		wantErr bool
	}{
		{name: "valid", request: MetricCreateRequest{CPU: 10, Memory: 20, Disk: 30}},
		{name: "negative cpu", request: MetricCreateRequest{CPU: -1, Memory: 20, Disk: 30}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}
