package handlers

import (
	mathbattle "mathbattle/models"

	tb "gopkg.in/tucnak/telebot.v2"
)

type StartRound struct {
	Handler
	Rounds mathbattle.RoundRepository
}

func (h *StartRound) Name() string {
	return h.Handler.Name
}

func (h *StartRound) Description() string {
	return h.Handler.Description
}

func (h *StartRound) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	res, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *StartRound) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, error) {
	return false, nil
}

func (h *StartRound) IsAdminOnly() bool {
	return true
}

func (h *StartRound) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	return -1, noResponse(), nil
}
