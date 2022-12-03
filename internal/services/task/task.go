package task

import (
	"context"
	"strings"

	"github.com/fr13n8/go-practice/internal/domain"
	"github.com/fr13n8/go-practice/internal/repository"
	"github.com/opentracing/opentracing-go"
)

type Task struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Task {
	return &Task{
		repo: repo,
	}
}

func (s *Task) Create(ctx context.Context, task domain.TaskCreate) (domain.Task, error) {
	span, jCtx := opentracing.StartSpanFromContext(ctx, "task.service.Create")
	defer span.Finish()

	item, err := s.repo.TaskDb.Create(jCtx, task)
	if err != nil {
		return domain.Task{}, err
	}
	err = s.repo.TaskRds.Created(jCtx, item)
	if err != nil {
		return domain.Task{}, err
	}

	err = s.repo.TaskRds.Set(jCtx, item, 0)
	if err != nil {
		return domain.Task{}, err
	}

	return item, nil
}

func (s *Task) Delete(ctx context.Context, id string) error {
	span, jCtx := opentracing.StartSpanFromContext(ctx, "task.service.Delete")
	defer span.Finish()

	return s.repo.TaskDb.Delete(jCtx, id)
}

func (s *Task) Update(ctx context.Context, task domain.TaskUpdate, id string) (domain.Task, error) {
	span, jCtx := opentracing.StartSpanFromContext(ctx, "task.service.Update")
	defer span.Finish()

	updTask := domain.Task{
		Name: task.Name,
		ID:   id,
	}
	newTask, err := s.repo.TaskDb.Update(jCtx, updTask)
	if err != nil {
		return domain.Task{}, err
	}

	err = s.repo.TaskRds.Set(jCtx, updTask, 0)
	if err != nil {
		return domain.Task{}, err
	}

	return newTask, nil
}

func (s *Task) Get(ctx context.Context, id string) (domain.Task, error) {
	span, jCtx := opentracing.StartSpanFromContext(ctx, "task.service.Get")
	defer span.Finish()

	task, err := s.repo.TaskRds.Get(jCtx, id)
	if err != nil {
		if !strings.Contains(err.Error(), "redis: nil") {
			return domain.Task{}, err
		}
	}
	if err == nil {
		return task, nil
	}

	task, err = s.repo.TaskDb.Get(jCtx, id)
	if err != nil {
		return domain.Task{}, err
	}

	err = s.repo.TaskRds.Set(jCtx, task, 0)
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (s *Task) GetAll(ctx context.Context) ([]domain.Task, error) {
	span, jCtx := opentracing.StartSpanFromContext(ctx, "task.service.GetAll")
	defer span.Finish()

	tasks, err := s.repo.TaskDb.GetAll(jCtx)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
