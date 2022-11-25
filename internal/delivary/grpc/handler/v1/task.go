package v1

import (
	"context"

	"github.com/fr13n8/go-practice/internal/domain"
	"github.com/fr13n8/go-practice/internal/services"
	pb "github.com/fr13n8/go-practice/pkg/grpc/v1"
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
	reqBody := domain.TaskCreate{}
	if req.GetName() == "" {
		return nil, nil
	}
	reqBody.Name = req.GetName()
	task, err := h.service.Create(reqBody)
	if err != nil {
		return nil, err
	}
	return &pb.CreateTaskResponse{
		Name:   task.Name,
		Id:     task.ID,
		Status: task.Status,
	}, nil
}

func (h *TaskHandler) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	id := req.GetId()
	task, err := h.service.Get(id)
	if err != nil {
		return nil, err
	}

	return &pb.GetTaskResponse{
		Name:   task.Name,
		Id:     task.ID,
		Status: task.Status,
	}, nil
}

func (h *TaskHandler) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	id := req.GetId()
	_, err := h.service.Get(id)
	if err != nil {
		return nil, err
	}

	reqBody := domain.TaskUpdate{}
	if req.GetName() == "" {
		return nil, nil
	}
	reqBody.Name = req.GetName()

	task, err := h.service.Update(reqBody, id)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateTaskResponse{
		Name:   task.Name,
		Id:     task.ID,
		Status: task.Status,
	}, nil
}

func (h *TaskHandler) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	id := req.GetId()
	_, err := h.service.Get(id)
	if err != nil {
		return nil, err
	}

	err = h.service.Delete(id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteTaskResponse{
		Id: id,
	}, nil
}

func (h *TaskHandler) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	tasks, err := h.service.GetAll()
	if err != nil {
		return nil, err
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
