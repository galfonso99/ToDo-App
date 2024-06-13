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

func (s *TaskStore) CreateTask(task *store.Task) (uint, error) {
    result := s.db.Create(task)
    err := result.Error
    return task.ID, err
}

func (s *TaskStore) GetTask(ID uint, sessionID string) (*store.Task, error) {
	var task store.Task
	err := s.db.Where("id = ? AND session_id = ?", ID, sessionID).First(&task).Error

	if err != nil {
		return nil, err
	}
	return &task, err
}

func (s *TaskStore) DeleteTask(ID uint) error {
	err := s.db.Delete(&store.Task{}, ID).Error
	return err
}

func (s *TaskStore) EditTask(ID uint, sessionID string, desc string) error {
    err := s.db.Save(&store.Task{ID: ID, SessionID: sessionID, Description: desc}).Error
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

func (s *TaskStore) GetTasksFromSession(sessionID string) ([]store.Task, error) {
    var tasks  []store.Task
	result := s.db.Where("session_id = ?", sessionID).Find(&tasks)
    err := result.Error

	if err != nil {
		return nil, err
	}
	return tasks, err
}
