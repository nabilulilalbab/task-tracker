package tests

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/nabilulilalbab/welcomesite/models"
	"github.com/nabilulilalbab/welcomesite/repositories"
	"github.com/nabilulilalbab/welcomesite/services"
	mockRepo "github.com/nabilulilalbab/welcomesite/tests/mock"
)

var dummyJPEG = []byte{
	// Minimal valid JPEG file (1x1 px)
	0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10,
	0x4A, 0x46, 0x49, 0x46, 0x00, 0x01,
	0x01, 0x00, 0x00, 0x01, 0x00, 0x01,
	0x00, 0x00, 0xFF, 0xDB, 0x00, 0x43,
	0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xC0, 0x00, 0x0B, 0x08, 0x00,
	0x01, 0x00, 0x01, 0x01, 0x01, 0x11,
	0x00, 0xFF, 0xC4, 0x00, 0x14, 0x00,
	0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0xFF, 0xDA, 0x00,
	0x08, 0x01, 0x01, 0x00, 0x00, 0x3F,
	0x00, 0xD2, 0xCF, 0x20, 0xFF, 0xD9,
}

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}
	if err := db.AutoMigrate(&models.Task{}); err != nil {
		panic("failed to migrate Task model")
	}
	return db
}

// ✅ Membuat *multipart.FileHeader dengan isi dummy (simulasi upload file)
func createMultipartFileHeader(t *testing.T, fieldName, fileName string, content []byte) *multipart.FileHeader {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	require.NoError(t, err)

	_, err = part.Write(content)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = req.ParseMultipartForm(int64(len(body.Bytes())))
	require.NoError(t, err)

	fileHeaders := req.MultipartForm.File[fieldName]
	require.NotEmpty(t, fileHeaders)

	return fileHeaders[0]
}

func TestCreateTask(t *testing.T) {
	db := setupTestDB()
	repo := repositories.NewTaskRepository(db)
	service := services.NewTaskService(repo)

	tests := []struct {
		name        string
		taskInput   *models.Task
		coverFile   *multipart.FileHeader
		expectError bool
	}{
		{
			name: "success with cover",
			taskInput: &models.Task{
				Judul: "Tugas 1",
			},

			coverFile:   createMultipartFileHeader(t, "cover", "cover.jpg", dummyJPEG),
			expectError: false,
		},
		{
			name: "success without cover",
			taskInput: &models.Task{
				Judul: "Tugas 2",
			},
			coverFile:   nil,
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.CreateTask(tc.taskInput, tc.coverFile)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.taskInput.Judul, result.Judul)

				if tc.coverFile != nil {
					assert.NotEmpty(t, result.Cover)
					assert.Contains(t, result.Cover, "/static/uploads/tasks/")
				} else {
					assert.Empty(t, result.Cover)
				}
			}
		})
	}
}

func TestGetTaskByID(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		mockReturn  *models.Task
		mockError   error
		expectError bool
	}{
		{
			name: "success",
			id:   1,
			mockReturn: &models.Task{
				ID:    1,
				Judul: "Test Judul",
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "task not found",
			id:          2,
			mockReturn:  nil,
			mockError:   errors.New("not found"),
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mockRepo.MockRepository)
			taskService := services.NewTaskService(mockRepo)

			mockRepo.On("FindByID", tc.id).Return(tc.mockReturn, tc.mockError)

			result, err := taskService.GetTaskByID(tc.id)
			mockRepo.AssertExpectations(t)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockReturn, result)
			}
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	tests := []struct {
		name        string
		mockReturn  []models.Task
		mockError   error
		expectError bool
	}{
		{
			name: "success",
			mockReturn: []models.Task{
				{ID: 1, Judul: "Task 1"},
				{ID: 2, Judul: "Task 2"},
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "error fetching tasks",
			mockReturn:  nil,
			mockError:   errors.New("DB error"),
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mockRepo.MockRepository)
			taskService := services.NewTaskService(mockRepo)

			mockRepo.On("FindAll").Return(tc.mockReturn, tc.mockError)

			result, err := taskService.GetAllTasks()
			mockRepo.AssertExpectations(t)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockReturn, result)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name        string
		input       *models.Task
		mockReturn  *models.Task
		mockError   error
		expectError bool
	}{
		{
			name: "success",
			input: &models.Task{
				ID:    1,
				Judul: "Test Judul Update",
			},
			mockReturn: &models.Task{
				ID:    1,
				Judul: "Test Judul Update",
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name: "error update task",
			input: &models.Task{
				ID:    1,
				Judul: "Test Judul Update",
			},
			mockReturn:  nil,
			mockError:   errors.New("failed update"),
			expectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mockRepo.MockRepository)
			taskService := services.NewTaskService(mockRepo)

			mockRepo.On("Update", tc.input).Return(tc.mockReturn, tc.mockError)

			result, err := taskService.UpdateTask(tc.input)
			mockRepo.AssertExpectations(t)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockReturn, result)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		mockReturn  error
		expectError bool
	}{
		{
			name:        "success",
			id:          1,
			mockReturn:  nil,
			expectError: false,
		},
		{
			name:        "failed delete task",
			id:          1,
			mockReturn:  errors.New("failed delete"),
			expectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mockRepo.MockRepository)
			taskService := services.NewTaskService(mockRepo)

			mockRepo.On("Delete", tc.id).Return(tc.mockReturn)

			err := taskService.DeleteTask(tc.id)
			mockRepo.AssertExpectations(t)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
