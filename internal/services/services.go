package services

import (
	"context"

	"github.com/fr13n8/go-practice/internal/domain"
	"github.com/fr13n8/go-practice/internal/repository"
	"github.com/fr13n8/go-practice/internal/services/task"
)

type Task interface {
	Create(ctx context.Context, task domain.TaskCreate) (domain.Task, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, task domain.TaskUpdate, id string) (domain.Task, error)
	Get(ctx context.Context, id string) (domain.Task, error)
	GetAll(ctx context.Context) ([]domain.Task, error)
}

type Services struct {
	Task Task
}

func NewServices(repos *repository.Repository) *Services {
	return &Services{
		Task: task.NewService(repos),
	}
}
