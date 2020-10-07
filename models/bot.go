package models

import (
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
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

type TelegramUserContext struct {
	Variables      map[string]string
	ChatID         int64
	CurrentStep    int
	CurrentCommand string
	Bot            *tb.Bot
}

func (c *TelegramUserContext) Recipient() string {
	return strconv.FormatInt(c.ChatID, 10)
}

func (c *TelegramUserContext) SendText(msg string) error {
	_, err := c.Bot.Send(c, msg)
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
