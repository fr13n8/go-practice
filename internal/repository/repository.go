package repository

import (
	"context"

	"github.com/fr13n8/go-practice/internal/domain"
	postgres_repo "github.com/fr13n8/go-practice/internal/repository/task/postgres"
	redis_repo "github.com/fr13n8/go-practice/internal/repository/task/redis"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type TaskDb interface {
	Create(ctx context.Context, task domain.TaskCreate) (domain.Task, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, task domain.Task) (domain.Task, error)
	Get(ctx context.Context, id string) (domain.Task, error)
	GetAll(ctx context.Context) ([]domain.Task, error)
}

type TaskRds interface {
	Set(ctx context.Context, task domain.Task, expire int) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.Task, error)
	Created(ctx context.Context, task domain.Task) error
}

type Repository struct {
	TaskDb
	TaskRds
}

func NewRepository(db *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{
		TaskDb:  postgres_repo.NewTask(db),
		TaskRds: redis_repo.NewTask(redis),
	}
}
