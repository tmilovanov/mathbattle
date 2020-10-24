package handlers

import (
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Stat struct {
	Handler
	AllCommands []mathbattle.TelegramCommandHandler
	Replier     mreplier.Replier
}

func (h *Stat) Name() string {
	return h.Handler.Name
}

func (h *Stat) Description() string {
	return h.Handler.Description
}

func (h *Stat) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	res, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *Stat) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, error) {
	return true, nil
}

func (h *Stat) IsAdminOnly() bool {
	return true
}

func (h *Stat) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	return -1, noResponse(), nil
}
