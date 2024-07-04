package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/popeskul/special-status-check/internal/config"
	"github.com/popeskul/special-status-check/internal/handlers"
	"github.com/popeskul/special-status-check/internal/health"
	"github.com/popeskul/special-status-check/internal/middleware"
	"github.com/popeskul/special-status-check/internal/server"
	"github.com/popeskul/special-status-check/internal/services"
)

type App struct {
	server *server.Server
	cfg    *config.Config
}

func NewApp(cfg *config.Config) (*App, error) {
	service, err := services.NewService()
	if err != nil {
		log.Fatalf("Failed to initialize service: %v", err)
	}

	router := http.NewServeMux()
	handler := handlers.NewHandler(service)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	handlerWithLogger := middleware.LoggingMiddleware(handler.InitRoutes(router))

	srv := server.NewServer(&http.Server{
		Addr:         addr,
		Handler:      handlerWithLogger,
		ReadTimeout:  cfg.Server.Timeouts.Read,
		WriteTimeout: cfg.Server.Timeouts.Write,
	})

	return &App{
		server: srv,
		cfg:    cfg,
	}, nil
}

func (a *App) Run() {
	go startServer(a.server)
	go health.StartHealthCheckServer(a.cfg.Server.HealthCheckPort)
	waitForShutdown(a.server)
}

func startServer(srv *server.Server) {
	log.Printf("Starting api on %s", srv.Addr())
	if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server stopped with error: %v", err)
	}
}

func waitForShutdown(srv *server.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	log.Println("Shutting down api...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %+v", err)
	}

	log.Println("Server gracefully stopped")
}
