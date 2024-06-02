package store

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Task struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Description    string `json:"Description"`
}

type Session struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	SessionID string `json:"session_id"`
	UserID    uint   `json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"user"`
}

type UserStore interface {
	CreateUser(description string, password string) error
	GetUser(email string) (*User, error)
}

type TaskStore interface {
	CreateTask(description string) error
	GetTask(description string) (*Task, error)
}

type SessionStore interface {
	CreateSession(session *Session) (*Session, error)
	GetUserFromSession(sessionID string, userID string) (*User, error)
}
