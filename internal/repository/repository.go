package repository

import (
	"github.com/fr13n8/go-practice/internal/domain"
	sqlite_repo "github.com/fr13n8/go-practice/internal/repository/task/postgres"
	redis_repo "github.com/fr13n8/go-practice/internal/repository/task/redis"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type TaskDb interface {
	Create(task domain.TaskCreate) (domain.Task, error)
	Delete(id string) error
	Update(task domain.Task) (domain.Task, error)
	Get(id string) (domain.Task, error)
	GetAll() ([]domain.Task, error)
}

type TaskRds interface {
	Set(task domain.Task, expire int) error
	Delete(id string) error
	Get(id string) (domain.Task, error)
	Created(task domain.Task) error
}

type Repository struct {
	TaskDb
	TaskRds
}

func NewRepository(db *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{
		TaskDb:  sqlite_repo.NewTask(db),
		TaskRds: redis_repo.NewTask(redis),
	}
}
