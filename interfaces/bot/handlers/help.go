package handlers

import (
	mreplier "mathbattle/application"
	"mathbattle/infrastructure"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Help struct {
	Handler
	Replier mreplier.Replier
}

func (h *Help) Name() string {
	return h.Handler.Name
}

func (h *Help) Description() string {
	return h.Handler.Description
}

func (h *Help) IsShowInHelp(ctx infrastructure.TelegramUserContext) bool {
	return true
}

func (h *Help) IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error) {
	return true, "", nil
}

func (h *Help) IsAdminOnly() bool {
	return false
}

func (h *Help) Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	return -1, NewResps(h.Replier.GetHelpMessages()...), nil
}
