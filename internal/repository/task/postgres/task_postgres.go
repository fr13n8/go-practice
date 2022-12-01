package sqlite_repo

import (
	"github.com/fr13n8/go-practice/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewTask(db *gorm.DB) *Repo {
	return &Repo{db}
}

func (e *Repo) Create(task domain.TaskCreate) (domain.Task, error) {
	newTask := domain.Task{
		Name:   task.Name,
		ID:     uuid.NewString(),
		Status: true,
	}
	if err := e.db.Table("tasks").Create(&newTask).Error; err != nil {
		return domain.Task{}, err
	}

	return newTask, nil
}

func (e *Repo) Delete(id string) error {
	var task domain.Task
	if err := e.db.Table("tasks").Where("id = ?", id).First(&task).Error; err != nil {
		return err
	}

	if err := e.db.Table("tasks").Delete(&task).Error; err != nil {
		return err
	}

	return nil
}

func (e *Repo) Update(task domain.Task) (domain.Task, error) {
	var newTask domain.Task
	if err := e.db.Table("tasks").Where("id = ?", task.ID).First(&task).Error; err != nil {
		return domain.Task{}, err
	}

	newTask.Name = task.Name
	newTask.Status = task.Status
	if err := e.db.Table("tasks").Save(&newTask).Error; err != nil {
		return domain.Task{}, err
	}

	return newTask, nil
}

func (e *Repo) Get(id string) (domain.Task, error) {
	var task domain.Task
	if err := e.db.Table("tasks").Where("id = ?", id).First(&task).Error; err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (e *Repo) GetAll() ([]domain.Task, error) {
	var tasks []domain.Task
	if err := e.db.Table("tasks").Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}
