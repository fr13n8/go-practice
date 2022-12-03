package v1

import (
	"context"
	"errors"

	"github.com/fr13n8/go-practice/internal/domain"
	"github.com/fr13n8/go-practice/internal/services"
	pb "github.com/fr13n8/go-practice/pkg/grpc/v1"
	"github.com/opentracing/opentracing-go"
)

type TaskHandler struct {
	pb.UnimplementedTaskServiceServer

	service services.Task
}

func NewTaskHandler(svc services.Task) *TaskHandler {
	return &TaskHandler{
		service: svc,
	}
}

func (h *TaskHandler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	span, jCtx := opentracing.StartSpanFromContext(ctx, "task.Create")
	defer span.Finish()
	reqBody := domain.TaskCreate{}
	if req.GetName() == "" {
		return nil, errors.New("name is required")
	}
	reqBody.Name = req.GetName()
	task, err := h.service.Create(jCtx, reqBody)
	if err != nil {
		return nil, errors.New("something went wrong")
	}
	return &pb.CreateTaskResponse{
		Name:   task.Name,
		Id:     task.ID,
		Status: task.Status,
	}, nil
}

func (h *TaskHandler) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	span, jCtx := opentracing.StartSpanFromContext(ctx, "task.Get")
	defer span.Finish()
	id := req.GetId()
	task, err := h.service.Get(jCtx, id)
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	return &pb.GetTaskResponse{
		Name:   task.Name,
		Id:     task.ID,
		Status: task.Status,
	}, nil
}

func (h *TaskHandler) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	span, jCtx := opentracing.StartSpanFromContext(ctx, "task.Update")
	defer span.Finish()
	id := req.GetId()
	_, err := h.service.Get(jCtx, id)
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	reqBody := domain.TaskUpdate{}
	if req.GetName() == "" {
		return nil, errors.New("name is required")
	}
	reqBody.Name = req.GetName()

	task, err := h.service.Update(jCtx, reqBody, id)
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	return &pb.UpdateTaskResponse{
		Name:   task.Name,
		Id:     task.ID,
		Status: task.Status,
	}, nil
}

func (h *TaskHandler) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	span, jCtx := opentracing.StartSpanFromContext(ctx, "task.Delete")
	defer span.Finish()
	id := req.GetId()
	_, err := h.service.Get(jCtx, id)
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	err = h.service.Delete(jCtx, id)
	if err != nil {
		return nil, errors.New("something went wrong")
	}
	return &pb.DeleteTaskResponse{
		Id: id,
	}, nil
}

func (h *TaskHandler) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	span, jCtx := opentracing.StartSpanFromContext(ctx, "handler.ListTasks")
	defer span.Finish()
	tasks, err := h.service.GetAll(jCtx)
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	var tasksResponse []*pb.Task
	for _, task := range tasks {
		tasksResponse = append(tasksResponse, &pb.Task{
			Name:   task.Name,
			Id:     task.ID,
			Status: task.Status,
		})
	}

	return &pb.ListTasksResponse{
		Tasks: tasksResponse,
	}, nil
}
