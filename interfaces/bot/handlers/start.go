package handlers

import (
	mreplier "mathbattle/application"
	"mathbattle/infrastructure"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Start struct {
	Handler
	Replier mreplier.Replier
}

func (h *Start) Name() string {
	return h.Handler.Name
}

func (h *Start) Description() string {
	return h.Handler.Description
}

func (h *Start) IsShowInHelp(ctx infrastructure.TelegramUserContext) bool {
	return false
}

func (h *Start) IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error) {
	return true, "", nil
}

func (h *Start) IsAdminOnly() bool {
	return false
}

func (h *Start) Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	return -1, OneTextResp(h.Replier.GetStartMessage()), nil
}
