package services

import (
	"github.com/fr13n8/go-practice/internal/domain"
	"github.com/fr13n8/go-practice/internal/repository"
	"github.com/fr13n8/go-practice/internal/services/task"
)

type Task interface {
	Create(task domain.TaskCreate) (domain.Task, error)
	Delete(id string) error
	Update(task domain.TaskUpdate, id string) (domain.Task, error)
	Get(id string) (domain.Task, error)
	GetAll() ([]domain.Task, error)
}

type Services struct {
	Task
}

func NewServices(repos *repository.Repository) *Services {
	return &Services{
		Task: task.NewService(repos),
	}
}
