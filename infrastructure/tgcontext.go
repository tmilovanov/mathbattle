package infrastructure

import (
	"strconv"

	"mathbattle/models/mathbattle"
)

type TelegramUserContext struct {
	User           mathbattle.User
	Variables      map[string]ContextVariable
	CurrentStep    int
	CurrentCommand string
}

type TelegramUserData struct {
	ChatID    int64
	FirstName string
	LastName  string
	Username  string
}

type TelegramContextRepository interface {
	GetByUserData(userData TelegramUserData) (TelegramUserContext, error)
	Update(chatID int64, ctx TelegramUserContext) error
}

func NewTelegramUserContext(userData TelegramUserData) TelegramUserContext {
	return TelegramUserContext{
		User: mathbattle.User{
			TelegramID:        userData.ChatID,
			TelegramFirstName: userData.FirstName,
			TelegramLastName:  userData.LastName,
			TelegramUsername:  userData.Username,
			IsAdmin:           false,
		},
		Variables:      make(map[string]ContextVariable),
		CurrentStep:    0,
		CurrentCommand: "",
	}
}

func NewTelegramUserContextByChatID(chatID int64) TelegramUserContext {
	return NewTelegramUserContext(TelegramUserData{
		ChatID:   chatID,
		Username: "",
	})
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
