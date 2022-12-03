package handler

import (
	"net/http"

	v1 "github.com/fr13n8/go-practice/internal/delivary/grpc/handler/v1"
	"github.com/fr13n8/go-practice/internal/interceptors"
	"github.com/fr13n8/go-practice/internal/services"
	pb "github.com/fr13n8/go-practice/pkg/grpc/v1/gen"
	metric "github.com/fr13n8/go-practice/pkg/metrics"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

type Handler struct {
	services    *services.Services
	TaskHandler pb.TaskServiceServer
}

func NewHandler(svcs *services.Services) *Handler {
	return &Handler{
		services:    svcs,
		TaskHandler: v1.NewTaskHandler(svcs.Task),
	}
}

func (h *Handler) Init(srv *grpc.Server) {
	metrics, err := metric.CreateMetricsGrpc("task_service_grpc")
	if err != nil {
		panic(err)
	}
	interceptors.NewInterceptorManager(metrics)
	grpc_prometheus.Register(srv)
	http.Handle("/metrics", promhttp.Handler())
	pb.RegisterTaskServiceServer(srv, h.TaskHandler)
}
