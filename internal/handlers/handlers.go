package handlers

import (
	"math/rand"
	"net/http"
	"time"

	_ "github.com/popeskul/special-status-check/docs"
	"github.com/popeskul/special-status-check/internal/metrics"
	"github.com/popeskul/special-status-check/internal/services"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//go:generate mockery --name HandlersInterface --output ./mocks --filename handlers.go
type HandlersInterface interface {
	InitRoutes(router *http.ServeMux) http.Handler
	GetInternal(w http.ResponseWriter, r *http.Request)
}

var (
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Handlers struct {
	service services.ServiceInterface
}

// NewHandler creates a new Handlers with the necessary dependencies.
func NewHandler(service services.ServiceInterface) *Handlers {
	return &Handlers{
		service: service,
	}
}

func (h *Handlers) InitRoutes(router *http.ServeMux) http.Handler {
	router.HandleFunc("/success", metrics.MeasureDuration(h.GetSuccess))
	router.HandleFunc("/internal", metrics.MeasureDuration(h.GetInternal))
	router.HandleFunc("/random", metrics.MeasureDuration(h.PostRandom))

	router.HandleFunc("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	return router
}

// GetInternal godoc
// @Summary Returns 500 Internal Server Error
// @Description Always returns 500 Internal Server Error
// @Tags status
// @Produce json
// @Success 500 {string} string "Internal Server Error"
// @Router /internal [get]
func (h *Handlers) GetInternal(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

// PostRandom godoc
// @Summary Returns random status
// @Description Returns 200 OK, 500 Internal Server Error, or panics with 20% probability
// @Tags status
// @Produce json
// @Success 200 {string} string "OK"
// @Failure 500 {string} string "Internal Server Error"
// @Router /random [post]
func (h *Handlers) PostRandom(w http.ResponseWriter, r *http.Request) {
	if rng.Intn(100) < 20 {
		panic("random panic")
	}

	if rng.Intn(2) == 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetSuccess godoc
// @Summary Returns 200 OK
// @Description Always returns 200 OK
// @Tags status
// @Produce json
// @Success 200 {string} string "OK"
// @Router /success [get]
func (h *Handlers) GetSuccess(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
