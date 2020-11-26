package models

import (
	"strconv"
)

type TelegramUser struct {
	ID      string
	ChatID  int64
	IsAdmin bool
}

type TelegramUserRepository interface {
	Store(user TelegramUser) (TelegramUser, error)
	GetByID(ID string) (TelegramUser, error)
	GetByTelegramID(ID int64) (TelegramUser, error)
	GetOrCreateByTelegramID(ID int64) (TelegramUser, error)
}

type TelegramUserContext struct {
	User           TelegramUser
	Variables      map[string]ContextVariable
	CurrentStep    int
	CurrentCommand string
}

type TelegramContextRepository interface {
	GetByTelegramID(chatID int64) (TelegramUserContext, error)
	Update(chatID int64, ctx TelegramUserContext) error
}

func NewTelegramUserContext(chatID int64) TelegramUserContext {
	return TelegramUserContext{
		User: TelegramUser{
			ChatID:  chatID,
			IsAdmin: false,
		},
		Variables:      make(map[string]ContextVariable),
		CurrentStep:    0,
		CurrentCommand: "",
	}
}

type ContextVariable struct {
	impl string
}

func NewContextVariableInt(init int) ContextVariable {
	return ContextVariable{strconv.Itoa(init)}
}

func NewContextVariableStr(init string) ContextVariable {
	return ContextVariable{init}
}

func (v ContextVariable) AsInt() (int, error) {
	return strconv.Atoi(v.impl)
}

func (v ContextVariable) AsString() string {
	return v.impl
}
