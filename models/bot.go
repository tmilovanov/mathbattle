package models

import (
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

type CommandStep string

var (
	StepStart CommandStep = "StepStart"
	StepSame  CommandStep = "StepSame"
	StepNext  CommandStep = "StepNext"
)

type TelegramResponse struct {
	Text     string
	Keyboard *tb.ReplyMarkup
}

func NewResp(messageText string) TelegramResponse {
	return TelegramResponse{messageText, &tb.ReplyMarkup{
		ReplyKeyboardRemove: true,
	}}
}

func NewRespWithKeyboard(messageText string, buttonTexts ...string) TelegramResponse {
	keyboard := &tb.ReplyMarkup{
		ResizeReplyKeyboard: true,
	}

	buttons := []tb.Btn{}
	for _, txt := range buttonTexts {
		buttons = append(buttons, keyboard.Text(txt))
	}

	keyboard.Reply(keyboard.Row(buttons...))

	return TelegramResponse{messageText, keyboard}
}

type TelegramCommandHandler interface {
	Name() string
	Description() string
	IsShowInHelp(ctx TelegramUserContext) bool
	Handle(ctx TelegramUserContext, m *tb.Message) (int, TelegramResponse, error)
}

type TelegramContextRepository interface {
	GetByTelegramID(chatID int64) (TelegramUserContext, error)
	Update(chatID int64, ctx TelegramUserContext) error
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

type TelegramUserContext struct {
	ChatID         int64
	Variables      map[string]ContextVariable
	CurrentStep    int
	CurrentCommand string
}

func NewTelegramUserContext(chatID int64) TelegramUserContext {
	return TelegramUserContext{chatID, make(map[string]ContextVariable), 0, ""}
}

func (c *TelegramUserContext) Recipient() string {
	//return strconv.FormatInt(c.ChatID, 10)
	return ""
}

func FilterCommandsToShow(allCommands []TelegramCommandHandler, ctx TelegramUserContext) []TelegramCommandHandler {
	result := []TelegramCommandHandler{}

	for _, cmd := range allCommands {
		if cmd.IsShowInHelp(ctx) {
			result = append(result, cmd)
		}
	}

	return result
}

func fillPhotoReader(msg *tb.Message) error {
	return nil
}
