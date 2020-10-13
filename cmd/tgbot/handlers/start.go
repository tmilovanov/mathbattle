package handlers

import (
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Start struct {
	Handler
	AllCommands []mathbattle.TelegramCommandHandler
	Replier     mreplier.Replier
}

func (h *Start) Name() string {
	return h.Handler.Name
}

func (h *Start) Description() string {
	return h.Handler.Description
}

func (h *Start) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	return false
}

func (h *Start) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, error) {
	return true, nil
}

func (h *Start) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, mathbattle.TelegramResponse, error) {
	msg := h.Replier.GetStartMessage()
	return -1, mathbattle.NewResp(msg), nil
}
