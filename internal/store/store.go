package store


type User struct {
	Email    string `json:"email"`
	Password string `json:"-"`
	ID       uint   `gorm:"primaryKey" json:"id"`
}

type Task struct {
	Description string `json:"Description"`
	Session     Session
	SessionID   string
	ID          uint   `gorm:"primaryKey" json:"id"`
}

type Session struct {
    ID string `gorm:"primaryKey" json:"session_id"`
}

type UserStore interface {
	CreateUser(email string, password string) error
	GetUser(email string) (*User, error)
}

type TaskStore interface {
	CreateTask(task *Task) (uint, error)
	GetTask(ID uint, sessionID string) (*Task, error)
    DeleteTask(ID uint) error
    EditTask(ID uint, sessionID string, desc string) error
	GetAllTasks() ([]Task, error)
	GetTasksFromSession(sessionID string) ([]Task, error)
}

type SessionStore interface {
	CreateSession() (*Session, error)
    // GetTasksFromSession(sessionID string) ([]Task, error)
}
