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

type TelegramResponse string

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

func (c *TelegramUserContext) SendText(msg string) error {
	//_, err := c.Bot.Send(c, msg, &tb.ReplyMarkup{
	//ReplyKeyboardRemove: true,
	//})
	//return err

	return nil
}

func (c *TelegramUserContext) SendMessageWithKeyboard(messageText string, buttonTexts ...string) error {
	//keyboard := &tb.ReplyMarkup{
	//ResizeReplyKeyboard: true,
	//}

	//buttons := []tb.Btn{}
	//for _, txt := range buttonTexts {
	//buttons = append(buttons, keyboard.Text(txt))
	//}

	//keyboard.Reply(keyboard.Row(buttons...))

	//_, err := c.Bot.Send(c, messageText, keyboard)
	//return err

	return nil
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
