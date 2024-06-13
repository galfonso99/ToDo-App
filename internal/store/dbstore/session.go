package dbstore

import (
	"goth/internal/store"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionStore struct {
	db *gorm.DB
}

type NewSessionStoreParams struct {
	DB *gorm.DB
}

func NewSessionStore(params NewSessionStoreParams) *SessionStore {
	return &SessionStore{
		db: params.DB,
	}
}

func (s *SessionStore) CreateSession() (*store.Session, error) {

	sessionID := uuid.New().String()
    session := store.Session {
        ID: sessionID,
    }

	result := s.db.Create(&session)

	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

