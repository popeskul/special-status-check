package app

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/popeskul/special-status-check/internal/config"
	"github.com/popeskul/special-status-check/internal/handlers"
	"github.com/popeskul/special-status-check/internal/middleware"
	"github.com/popeskul/special-status-check/internal/server"
	"github.com/popeskul/special-status-check/internal/services"
)

// Mock handler to be used in tests
func mockHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestNewApp(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
			Timeouts: config.Timeouts{
				Read:  5 * time.Second,
				Write: 10 * time.Second,
				Idle:  15 * time.Second,
			},
			HealthCheckPort: 9090,
		},
	}

	app, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	if app.server.Addr() != ":8080" {
		t.Errorf("Unexpected server address: got %v, want %v", app.server.Addr(), ":8080")
	}
}

func TestApp_Run(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
			Timeouts: config.Timeouts{
				Read:  5 * time.Second,
				Write: 10 * time.Second,
				Idle:  15 * time.Second,
			},
			HealthCheckPort: 9090,
		},
	}

	// Mock services and handlers
	service, _ := services.NewService()
	router := http.NewServeMux()
	handler := handlers.NewHandler(service)
	handler.InitRoutes(router)
	handlerWithLogger := middleware.LoggingMiddleware(http.HandlerFunc(mockHandler))

	app := &App{
		server: server.NewServer(&http.Server{
			Addr:         ":0", // Use dynamic port
			Handler:      handlerWithLogger,
			ReadTimeout:  cfg.Server.Timeouts.Read,
			WriteTimeout: cfg.Server.Timeouts.Write,
			IdleTimeout:  cfg.Server.Timeouts.Idle,
		}),
		cfg: cfg,
	}

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{"API server", "/", http.StatusOK},
		{"Health check server", "/healthz", http.StatusOK},
	}

	// Start the server in a separate goroutine
	server := httptest.NewServer(handlerWithLogger)
	defer server.Close()

	// Mock health check server
	healthServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer healthServer.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var url string
			if tt.name == "API server" {
				url = server.URL + tt.path
			} else {
				url = healthServer.URL + tt.path
			}

			resp, err := http.Get(url)
			if err != nil {
				t.Fatalf("Failed to perform GET request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Unexpected status code: got %v, want %v", resp.StatusCode, tt.expectedStatus)
			}
		})
	}

	// Simulate shutdown signal
	go func() {
		time.Sleep(100 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		err := p.Signal(os.Interrupt)
		if err != nil {
			t.Errorf("Failed to send signal: %v", err)
			return
		}
	}()

	// Wait for shutdown
	waitForShutdown(app.server)
}

func TestWaitForShutdown(t *testing.T) {
	srv := server.NewServer(&http.Server{
		Addr:    ":0",                          // Use dynamic port
		Handler: http.HandlerFunc(mockHandler), // Mock handler for testing
	})

	// Start the server in a separate goroutine
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("Server stopped with error: %v", err)
			return
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Simulate shutdown signal
	go func() {
		time.Sleep(100 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		err := p.Signal(os.Interrupt)
		if err != nil {
			t.Errorf("Failed to send signal: %v", err)
			return
		}
	}()

	waitForShutdown(srv)

	// Add delay to ensure server has shut down
	time.Sleep(100 * time.Millisecond)

	// Check if server is shut down by trying to connect
	resp, err := http.Get(srv.Addr())
	if err == nil {
		resp.Body.Close()
		t.Errorf("Expected server to be shut down, but it is still running")
	}
}
