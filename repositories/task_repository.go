package repositories

import (
	"gorm.io/gorm"

	"github.com/nabilulilalbab/welcomesite/models"
)

type TaskRepository interface {
	Create(task *models.Task) (*models.Task, error)
	FindByID(id uint) (*models.Task, error)
	FindAll() ([]models.Task, error)
	Update(task *models.Task) (*models.Task, error)
	Delete(id uint) error
	GetDB() *gorm.DB
	FindByIDWithTx(id uint, tx *gorm.DB) (*models.Task, error)
}

type TaskRepositoryImpl struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &TaskRepositoryImpl{db: db}
}

func (r *TaskRepositoryImpl) GetDB() *gorm.DB {
	return r.db
}

func (t *TaskRepositoryImpl) Create(task *models.Task) (*models.Task, error) {
	err := t.db.Create(task).Error
	return task, err
}

func (t *TaskRepositoryImpl) FindByID(id uint) (*models.Task, error) {
	var task models.Task
	err := t.db.First(&task, id).Error
	return &task, err
}

func (t *TaskRepositoryImpl) FindAll() ([]models.Task, error) {
	var task []models.Task
	err := t.db.Find(&task).Error
	return task, err
}

func (t *TaskRepositoryImpl) Update(task *models.Task) (*models.Task, error) {
	err := t.db.Save(task).Error
	return task, err
}

func (t *TaskRepositoryImpl) Delete(id uint) error {
	err := t.db.Delete(&models.Task{}, id).Error
	return err
}

func (t *TaskRepositoryImpl) FindByIDWithTx(id uint, tx *gorm.DB) (*models.Task, error) {
	var task models.Task
	if err := tx.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}
