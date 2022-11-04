package task

import (
	"github.com/fr13n8/go-practice/internal/domain"
	"github.com/fr13n8/go-practice/internal/repository"
	"strings"
)

type Task struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Task {
	return &Task{
		repo: repo,
	}
}

func (s *Task) Create(task domain.TaskCreate) (domain.Task, error) {
	item, err := s.repo.TaskDb.Create(task)
	if err != nil {
		return domain.Task{}, err
	}
	err = s.repo.TaskRds.Created(item)
	if err != nil {
		return domain.Task{}, err
	}

	return item, nil
}

func (s *Task) Delete(id string) error {
	return s.repo.TaskDb.Delete(id)
}

func (s *Task) Update(task domain.TaskUpdate, id string) (domain.Task, error) {
	updTask := domain.Task{
		Name: task.Name,
		ID:   id,
	}
	newTask, err := s.repo.TaskDb.Update(updTask)
	if err != nil {
		return domain.Task{}, err
	}

	return newTask, nil
}

func (s *Task) Get(id string) (domain.Task, error) {
	task, err := s.repo.TaskRds.Get(id)
	if err != nil {
		if !strings.Contains(err.Error(), "redis: nil") {
			return domain.Task{}, err
		}
	}
	if err == nil {
		return task, nil
	}

	task, err = s.repo.TaskDb.Get(id)
	if err != nil {
		return domain.Task{}, err
	}

	err = s.repo.TaskRds.Set(task, 0)
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (s *Task) GetAll() ([]domain.Task, error) {
	tasks, err := s.repo.TaskDb.GetAll()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
