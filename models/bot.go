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

type TelegramCommandHandler interface {
	Name() string
	Description() string
	IsShowInHelp(ctx TelegramUserContext) bool
	Handle(ctx TelegramUserContext, m *tb.Message) (int, error)
}

type TelegramContextRepository interface {
	GetByTelegramID(chatID int64, bot *tb.Bot) (TelegramUserContext, error)
	Update(ctx TelegramUserContext) error
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
	Variables      map[string]ContextVariable
	ChatID         int64
	CurrentStep    int
	CurrentCommand string
	Bot            *tb.Bot
}

func (c *TelegramUserContext) Recipient() string {
	return strconv.FormatInt(c.ChatID, 10)
}

func (c *TelegramUserContext) SendText(msg string) error {
	_, err := c.Bot.Send(c, msg, &tb.ReplyMarkup{
		ReplyKeyboardRemove: true,
	})
	return err
}

func (c *TelegramUserContext) SendMessageWithKeyboard(messageText string, buttonTexts ...string) error {
	keyboard := &tb.ReplyMarkup{
		ResizeReplyKeyboard: true,
	}

	buttons := []tb.Btn{}
	for _, txt := range buttonTexts {
		buttons = append(buttons, keyboard.Text(txt))
	}

	keyboard.Reply(keyboard.Row(buttons...))

	_, err := c.Bot.Send(c, messageText, keyboard)
	return err
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
