package health

import (
	"fmt"
	"log"
	"net/http"

	"github.com/popeskul/special-status-check/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartHealthCheckServer(port int) {
	metrics.InitMetrics()

	router := http.NewServeMux()
	router.HandleFunc("/healthz", metrics.MeasureDuration(getHealthz))
	router.HandleFunc("/readyz", metrics.MeasureDuration(getReadyz))
	router.HandleFunc("/metrics", getMetrics)

	log.Printf("Starting health check server on port %d\n", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
		log.Fatalf("Failed to start health check server: %v", err)
	}
}

func getHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getReadyz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getMetrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}
