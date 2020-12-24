package mathbattle

type User struct {
	ID      string
	ChatID  int64
	IsAdmin bool
}

type UserRepository interface {
	Store(user User) (User, error)
	GetByID(ID string) (User, error)
	GetByTelegramID(ID int64) (User, error)
	GetOrCreateByTelegramID(ID int64) (User, error)
}
