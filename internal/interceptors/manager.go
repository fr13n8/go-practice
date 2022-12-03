package interceptors

import (
	"context"
	"fmt"
	"net/http"
	"time"

	metric "github.com/fr13n8/go-practice/pkg/metrics"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

// InterceptorManager
type InterceptorManager struct {
	metr metric.Metrics
}

// InterceptorManager constructor
func NewInterceptorManager(metr metric.Metrics) *InterceptorManager {
	return &InterceptorManager{metr: metr}
}

func (im *InterceptorManager) Metrics(ctx *fiber.Ctx) error {
	start := time.Now()
	next := ctx.Next()
	status := http.StatusOK
	if ctx.Response().StatusCode() != 0 {
		status = ctx.Response().StatusCode()
	}

	im.metr.ObserveResponseTime(status, ctx.Method(), string(ctx.Context().Request.URI().Path()), time.Since(start).Seconds())
	im.metr.IncHits(status, ctx.Method(), string(ctx.Context().Request.URI().Path()))

	return next
}

func (im *InterceptorManager) GrpcMetrics(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	var status = http.StatusOK
	if err != nil {
		status = http.StatusInternalServerError
	}
	fmt.Println("info.FullMethod", info.FullMethod)
	im.metr.ObserveResponseTime(status, info.FullMethod, info.FullMethod, time.Since(start).Seconds())
	im.metr.IncHits(status, info.FullMethod, info.FullMethod)

	return resp, err
}
