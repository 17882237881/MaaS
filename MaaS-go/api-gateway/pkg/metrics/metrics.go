package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTPRequestDuration tracks HTTP request duration
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestTotal tracks total HTTP requests
	HTTPRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestSize tracks HTTP request size
	HTTPRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	)

	// HTTPResponseSize tracks HTTP response size
	HTTPResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	)

	// ActiveConnections tracks active connections
	ActiveConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_connections",
			Help: "Number of active HTTP connections",
		},
	)

	// ServiceUp indicates if service is up
	ServiceUp = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "service_up",
			Help: "Whether the service is up (1) or down (0)",
		},
	)

	// ServiceInfo provides service information
	ServiceInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_info",
			Help: "Service information",
		},
		[]string{"version", "environment"},
	)
)

func init() {
	// Register all metrics
	prometheus.MustRegister(HTTPRequestDuration)
	prometheus.MustRegister(HTTPRequestTotal)
	prometheus.MustRegister(HTTPRequestSize)
	prometheus.MustRegister(HTTPResponseSize)
	prometheus.MustRegister(ActiveConnections)
	prometheus.MustRegister(ServiceUp)
	prometheus.MustRegister(ServiceInfo)
}

// PrometheusMiddleware returns a Gin middleware that collects Prometheus metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip metrics endpoint itself
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		// Track active connections
		ActiveConnections.Inc()
		defer ActiveConnections.Dec()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		// Record metrics
		HTTPRequestDuration.WithLabelValues(method, path, status).Observe(duration)
		HTTPRequestTotal.WithLabelValues(method, path, status).Inc()
		HTTPRequestSize.WithLabelValues(method, path).Observe(float64(c.Request.ContentLength))
		HTTPResponseSize.WithLabelValues(method, path).Observe(float64(c.Writer.Size()))
	}
}

// Handler returns HTTP handler for Prometheus metrics
func Handler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

// SetServiceUp sets service_up metric
func SetServiceUp(up bool) {
	if up {
		ServiceUp.Set(1)
	} else {
		ServiceUp.Set(0)
	}
}

// SetServiceInfo sets service_info metric
func SetServiceInfo(version, environment string) {
	ServiceInfo.WithLabelValues(version, environment).Set(1)
}

// RecordCustomMetric records a custom counter metric
func RecordCustomMetric(name string, value float64, labels ...string) {
	// This is a placeholder for custom metrics
	// In production, you would define custom metrics and use them here
}
