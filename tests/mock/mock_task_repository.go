package mock

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/nabilulilalbab/welcomesite/models"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(task *models.Task) (*models.Task, error) {
	args := m.Called(task)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockRepository) FindByID(id uint) (*models.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockRepository) FindAll() ([]models.Task, error) {
	args := m.Called()
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockRepository) Update(task *models.Task) (*models.Task, error) {
	args := m.Called(task)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) GetDB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}
