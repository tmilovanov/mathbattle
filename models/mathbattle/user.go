package mathbattle

import "time"

type User struct {
	ID               string    `json:"id"`
	TelegramID       int64     `json:"telegram_id"`
	TelegramName     string    `json:"telegram_name"`
	IsAdmin          bool      `json:"is_admin"`
	RegistrationTime time.Time `json:"registration_time"`
}

type UserRepository interface {
	Store(user User) (User, error)
	GetAll() ([]User, error)
	GetByID(ID string) (User, error)
	GetByTelegramID(ID int64) (User, error)
	GetByTelegramName(name string) (User, error)
	Update(user User) error
}
