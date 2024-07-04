package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TestMeasureDuration(t *testing.T) {
	InitMetrics()

	tests := []struct {
		name           string
		method         string
		path           string
		handler        http.HandlerFunc
		expectedMethod string
		expectedPath   string
	}{
		{
			name:   "GET request",
			method: http.MethodGet,
			path:   "/test-get",
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(10 * time.Millisecond)
				w.WriteHeader(http.StatusOK)
			},
			expectedMethod: http.MethodGet,
			expectedPath:   "/test-get",
		},
		{
			name:   "POST request",
			method: http.MethodPost,
			path:   "/test-post",
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(20 * time.Millisecond)
				w.WriteHeader(http.StatusOK)
			},
			expectedMethod: http.MethodPost,
			expectedPath:   "/test-post",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Prometheus registry for this test
			reg := prometheus.NewRegistry()
			reg.MustRegister(HttpRequestDuration)

			handler := MeasureDuration(tt.handler)
			server := httptest.NewServer(http.HandlerFunc(handler))
			defer server.Close()

			req, err := http.NewRequest(tt.method, server.URL+tt.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			client := &http.Client{Timeout: 1 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to perform request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Unexpected status code: got %v, want %v", resp.StatusCode, http.StatusOK)
			}

			// Gather metrics and check the value
			mfs, err := reg.Gather()
			if err != nil {
				t.Fatalf("Failed to gather metrics: %v", err)
			}

			found := false
			for _, mf := range mfs {
				if mf.GetName() == "http_request_duration_seconds" {
					for _, m := range mf.GetMetric() {
						labels := m.GetLabel()
						if len(labels) == 2 && labels[0].GetValue() == tt.expectedMethod && labels[1].GetValue() == tt.expectedPath {
							found = true
							if m.GetHistogram().GetSampleCount() == 0 {
								t.Errorf("Expected metric with labels %v to be present and non-zero", prometheus.Labels{"method": tt.expectedMethod, "path": tt.expectedPath})
							}
						}
					}
				}
			}

			if !found {
				t.Errorf("Expected metric with labels %v to be present", prometheus.Labels{"method": tt.expectedMethod, "path": tt.expectedPath})
			}
		})
	}
}
