package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/nabilulilalbab/welcomesite/models"
	"github.com/nabilulilalbab/welcomesite/repositories"
	"github.com/nabilulilalbab/welcomesite/utils"
)

type TaskService interface {
	CreateTask(task *models.Task, coverFile *multipart.FileHeader) (*models.Task, error)
	GetTaskByID(id uint) (*models.Task, error)
	GetAllTasks() ([]models.Task, error)
	UpdateTask(id uint, task *models.Task, fileHeader *multipart.FileHeader) (*models.Task, error)
	DeleteTask(id uint) error
}

type taskServiceImpl struct {
	repo        repositories.TaskRepository
	uploadsPath string
}

func NewTaskService(repository repositories.TaskRepository, uploadsPath string) TaskService {
	return &taskServiceImpl{repo: repository, uploadsPath: uploadsPath}
}

func (s *taskServiceImpl) CreateTask(task *models.Task, coverFile *multipart.FileHeader) (*models.Task, error) {
	tx := s.repo.GetDB().Begin()
	fmt.Println("tx : ", *tx)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if err := tx.Create(task).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if coverFile != nil {
		uniqueFileName := "task_" + strconv.FormatUint(uint64(task.ID), 10) + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(coverFile.Filename)
		diskPath := filepath.Join(s.uploadsPath, uniqueFileName)
		src, err := coverFile.Open()
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		defer src.Close()
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, src); err != nil {
			tx.Rollback()
			return nil, err
		}
		err = utils.SaveResizedImage(bytes.NewReader(buf.Bytes()), filepath.Ext(coverFile.Filename), diskPath, 800)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		task.Cover = "/static/uploads/tasks/" + uniqueFileName
		if err := tx.Model(task).Update("cover", task.Cover).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskServiceImpl) UpdateTask(id uint, taskInput *models.Task, coverFile *multipart.FileHeader) (*models.Task, error) {
	tx := s.repo.GetDB().Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	existingTask, err := s.repo.FindByIDWithTx(id, tx)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("task with id %d not found", id)
	}
	if err := tx.Model(existingTask).Updates(taskInput).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if coverFile != nil {
		if existingTask.Cover != "" {
			oldPath := filepath.Join(".", existingTask.Cover)
			if err := os.Remove(oldPath); err != nil {
				fmt.Printf("Warning: could not remove old file %s: %v\n", oldPath, err)
			}
		}
		uniqueFileName := "task_" + strconv.FormatUint(uint64(existingTask.ID), 10) + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(coverFile.Filename)
		diskPath := filepath.Join(s.uploadsPath, uniqueFileName)
		src, err := coverFile.Open()
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		defer src.Close()
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, src); err != nil {
			tx.Rollback()
			return nil, err
		}
		err = utils.SaveResizedImage(bytes.NewReader(buf.Bytes()), filepath.Ext(coverFile.Filename), diskPath, 800)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		newCoverPath := "/static/uploads/tasks/" + uniqueFileName
		if err := tx.Model(existingTask).Update("cover", newCoverPath).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		existingTask.Cover = newCoverPath
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return existingTask, nil
}

func (s *taskServiceImpl) GetTaskByID(id uint) (*models.Task, error) {
	return s.repo.FindByID(id)
}

func (s *taskServiceImpl) GetAllTasks() ([]models.Task, error) {
	return s.repo.FindAll()
}

func (s *taskServiceImpl) DeleteTask(id uint) error {
	// 1. Ambil data tugas untuk mendapatkan path cover
	task, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("gagal menemukan tugas dengan ID %d: %w", id, err)
	}

	// 2. Hapus file cover jika ada
	if task.Cover != "" {
		// Path cover di database adalah URL, kita perlu mengubahnya menjadi path fisik
		// Contoh: /static/uploads/tasks/task_1_...png -> static/uploads/tasks/task_1_...png
		// Kita perlu menghapus '/static/' dari awal URL
		fullPath := filepath.Join(s.uploadsPath, filepath.Base(task.Cover))

		err := os.Remove(fullPath)
		if err != nil {
			// Log error tapi jangan kembalikan error, agar penghapusan database tetap berjalan
			// karena penghapusan file mungkin gagal karena file tidak ada, izin, dll.
			fmt.Printf("Peringatan: Gagal menghapus file cover %s: %v\n", fullPath, err)
		}
	}

	// 3. Hapus data dari database
	return s.repo.Delete(id)
}
