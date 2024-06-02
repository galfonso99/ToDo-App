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

func (s *TaskStore) CreateTask(task string) error {
	return s.db.Create(&store.Task{
		Description:    task,
	}).Error
}

func (s *TaskStore) GetTask(email string) (*store.Task, error) {

	var user store.Task
	err := s.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, err
}

