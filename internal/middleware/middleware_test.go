package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLoggingMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		target     string
		wantStatus int
	}{
		{
			name:       "GET request",
			method:     http.MethodGet,
			target:     "/test-get",
			wantStatus: http.StatusOK,
		},
		{
			name:       "POST request",
			method:     http.MethodPost,
			target:     "/test-post",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server with the middleware
			handler := LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			server := httptest.NewServer(handler)
			defer server.Close()

			// Create a request to the test server
			req, err := http.NewRequest(tt.method, server.URL+tt.target, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Capture the logs
			var logs strings.Builder
			log.SetOutput(&logs)
			defer func() {
				log.SetOutput(nil) // Reset log output
			}()

			// Perform the request
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to perform request: %v", err)
			}
			defer resp.Body.Close()

			// Check the response status
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Unexpected status code: got %v, want %v", resp.StatusCode, tt.wantStatus)
			}

			// Check the logs
			if !strings.Contains(logs.String(), "Started "+tt.method+" "+tt.target) {
				t.Errorf("Log does not contain start message")
			}
			if !strings.Contains(logs.String(), "Completed "+tt.method+" "+tt.target) {
				t.Errorf("Log does not contain completion message")
			}
		})
	}
}
