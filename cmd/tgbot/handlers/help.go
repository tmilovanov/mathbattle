package handlers

import (
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"

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

func (h *Help) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	return true
}

func (h *Help) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, string, error) {
	return true, "", nil
}

func (h *Help) IsAdminOnly() bool {
	return false
}

func (h *Help) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	msg := h.Replier.GetHelpMessage()
	return -1, mathbattle.OneTextResp(msg), nil
}
