package mathbattle

import "time"

type User struct {
	ID                string    `json:"id"`
	TelegramID        int64     `json:"telegram_id"`
	TelegramFirstName string    `json:"telegram_firstname"`
	TelegramLastName  string    `json:"telegram_lastname"`
	TelegramUsername  string    `json:"telegram_username"`
	IsAdmin           bool      `json:"is_admin"`
	RegistrationTime  time.Time `json:"registration_time"`
}

func (u *User) SetRegistrationTime(t time.Time) {
	u.RegistrationTime = t.Round(time.Second).UTC()
}

type UserRepository interface {
	Store(user User) (User, error)
	GetAll() ([]User, error)
	GetByID(ID string) (User, error)
	GetByTelegramID(ID int64) (User, error)
	GetByTelegramName(name string) (User, error)
	Update(user User) error
}
