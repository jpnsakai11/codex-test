package http

import (
	"embed"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/example/microservices-project/order-service/internal/domain"
	"github.com/example/microservices-project/order-service/internal/usecase"
)

//go:embed docs/openapi.yaml
var openAPISpec []byte

const swaggerUI = `<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>Order Service API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"/>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>window.ui = SwaggerUIBundle({url: '/swagger/openapi.yaml', dom_id: '#swagger-ui'});</script>
</body>
</html>`

type Handler struct {
	uc          *usecase.OrderUsecase
	logger      *zap.Logger
	requestDur  *prometheus.HistogramVec
	requestCnt  *prometheus.CounterVec
	serviceName string
}

func NewHandler(uc *usecase.OrderUsecase, logger *zap.Logger, serviceName string) *Handler {
	dur := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests.",
		Buckets: prometheus.DefBuckets,
	}, []string{"service", "method", "path", "status"})
	cnt := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests.",
	}, []string{"service", "method", "path", "status"})
	prometheus.MustRegister(dur, cnt)

	return &Handler{uc: uc, logger: logger, requestDur: dur, requestCnt: cnt, serviceName: serviceName}
}

func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(h.requestLogger)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	r.Get("/readyz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	r.Handle("/metrics", promhttp.Handler())
	r.Get("/swagger", h.swaggerUI)
	r.Get("/swagger/openapi.yaml", h.openapiSpec)
	r.Post("/orders", h.createOrder)
	r.Get("/orders/{id}", h.getOrder)

	return r
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var order domain.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		h.respondErr(w, http.StatusBadRequest, err)
		return
	}
	if err := h.uc.Create(r.Context(), &order); err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			h.respondErr(w, http.StatusNotFound, err)
			return
		}
		h.respondErr(w, http.StatusBadRequest, err)
		return
	}
	h.respondJSON(w, http.StatusCreated, order)
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.respondErr(w, http.StatusBadRequest, err)
		return
	}
	order, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		h.respondErr(w, http.StatusBadRequest, err)
		return
	}
	if order == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	h.respondJSON(w, http.StatusOK, order)
}

func (h *Handler) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(ww, r)
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(ww.status)
		path := r.URL.Path
		h.requestDur.WithLabelValues(h.serviceName, r.Method, path, status).Observe(duration)
		h.requestCnt.WithLabelValues(h.serviceName, r.Method, path, status).Inc()
		h.logger.Info("http_request",
			zap.String("method", r.Method),
			zap.String("path", path),
			zap.Int("status", ww.status),
			zap.Float64("duration_seconds", duration),
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (h *Handler) respondJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *Handler) respondErr(w http.ResponseWriter, code int, err error) {
	h.respondJSON(w, code, map[string]string{"error": err.Error()})
}

func (h *Handler) swaggerUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(swaggerUI))
}

func (h *Handler) openapiSpec(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(openAPISpec)
}
