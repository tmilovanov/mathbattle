package models

import tb "gopkg.in/tucnak/telebot.v2"

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
	IsCommandSuitable(ctx TelegramUserContext) (bool, error)
	Handle(ctx TelegramUserContext, m *tb.Message) (int, TelegramResponse, error)
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
