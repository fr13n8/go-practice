package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var latency = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "api",
		Name:       "latency_seconds",
		Help:       "latency distributions.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"method", "path"},
)

func RegisterPromethesMetrics() {
	prometheus.MustRegister(latency)
}

func RecrdRequestLatency(ctx *fiber.Ctx) error {
	start := time.Now()
	next := ctx.Next()
	elapsed := time.Since(start).Seconds()

	latency.WithLabelValues(ctx.Method(), string(ctx.Context().Request.URI().Path())).Observe(elapsed)

	return next
}
