package dbstore

import (
	"goth/internal/hash"
	"goth/internal/store"

	"gorm.io/gorm"
)

type TaskStore struct {
	db           *gorm.DB
}

type NewTaskStoreParams struct {
	DB           *gorm.DB
	PasswordHash hash.PasswordHash
}

func NewTaskStore(params NewTaskStoreParams) *TaskStore {
	return &TaskStore{
		db:           params.DB,
	}
}

func (s *TaskStore) CreateTask(description string) (uint, error) {
    task := store.Task {
        Description: description,
    }
    result := s.db.Create(&task)
    err := result.Error
    return task.ID, err
}

func (s *TaskStore) GetTask(ID uint) (*store.Task, error) {

	var task store.Task
	err := s.db.Where("ID = ?", ID).First(&task).Error

	if err != nil {
		return nil, err
	}
	return &task, err
}

func (s *TaskStore) DeleteTask(ID uint) error {
	err := s.db.Delete(&store.Task{}, ID).Error
	return err
}

func (s *TaskStore) EditTask(ID uint, desc string) error {
    err := s.db.Save(&store.Task{ID: ID, Description: desc}).Error
	return err
}

func (s *TaskStore) GetAllTasks() ([]store.Task, error) {
    var tasks  []store.Task
	result := s.db.Find(&tasks)
    err := result.Error

	if err != nil {
		return nil, err
	}
	return tasks, err
}

