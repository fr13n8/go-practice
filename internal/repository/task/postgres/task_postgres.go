package sqlite_repo

import (
	"context"

	"github.com/fr13n8/go-practice/internal/domain"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewTask(db *gorm.DB) *Repo {
	return &Repo{db}
}

func (e *Repo) Create(ctx context.Context, task domain.TaskCreate) (domain.Task, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "task.repository.Create")
	defer span.Finish()

	newTask := domain.Task{
		Name:   task.Name,
		ID:     uuid.NewString(),
		Status: true,
	}
	if err := e.db.Table("tasks").Create(&newTask).Error; err != nil {
		ext.LogError(span, err)
		return domain.Task{}, err
	}

	return newTask, nil
}

func (e *Repo) Delete(ctx context.Context, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "task.repository.Delete")
	defer span.Finish()

	var task domain.Task
	if err := e.db.Table("tasks").Where("id = ?", id).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return gorm.ErrRecordNotFound
		}
		ext.LogError(span, err)
		return err
	}

	if err := e.db.Table("tasks").Delete(&task).Error; err != nil {
		ext.LogError(span, err)
		return err
	}

	return nil
}

func (e *Repo) Update(ctx context.Context, task domain.Task) (domain.Task, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "task.repository.Update")
	defer span.Finish()

	var newTask domain.Task
	if err := e.db.Table("tasks").Where("id = ?", task.ID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Task{}, gorm.ErrRecordNotFound
		}
		ext.LogError(span, err)
		return domain.Task{}, err
	}

	newTask.ID = task.ID
	newTask.Name = task.Name
	newTask.Status = task.Status
	if err := e.db.Table("tasks").Save(&newTask).Error; err != nil {
		ext.LogError(span, err)
		return domain.Task{}, err
	}

	return newTask, nil
}

func (e *Repo) Get(ctx context.Context, id string) (domain.Task, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "task.repository.Get")
	defer span.Finish()

	var task domain.Task
	if err := e.db.Table("tasks").Where("id = ?", id).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Task{}, gorm.ErrRecordNotFound
		}
		ext.LogError(span, err)
		return domain.Task{}, err
	}

	return task, nil
}

func (e *Repo) GetAll(ctx context.Context) ([]domain.Task, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "task.repository.GetAll")
	defer span.Finish()

	var tasks []domain.Task
	if err := e.db.Table("tasks").Find(&tasks).Error; err != nil {
		ext.LogError(span, err)
		return nil, err
	}

	return tasks, nil
}
