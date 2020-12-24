package handlers

import (
	"mathbattle/application"
	"mathbattle/infrastructure"

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
	IsShowInHelp(ctx infrastructure.TelegramUserContext) bool
	IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error)
	IsAdminOnly() bool
	Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error)
}

func FilterCommandsToShow(allCommands []TelegramCommandHandler, ctx infrastructure.TelegramUserContext) []application.TelegramCommandHelp {
	result := []application.TelegramCommandHelp{}

	for _, cmd := range allCommands {
		if !ctx.User.IsAdmin && cmd.IsAdminOnly() {
			continue
		}

		if cmd.IsShowInHelp(ctx) {
			result = append(result, application.TelegramCommandHelp{
				Name: cmd.Name(),
				Desc: cmd.Description(),
			})
		}
	}

	return result
}
