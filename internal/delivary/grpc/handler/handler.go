package handler

import (
	pb "github.com/fr13n8/go-practice/pkg/grpc/v1"

	v1 "github.com/fr13n8/go-practice/internal/delivary/grpc/handler/v1"
	"github.com/fr13n8/go-practice/internal/services"
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
	pb.RegisterTaskServiceServer(srv, h.TaskHandler)
}
