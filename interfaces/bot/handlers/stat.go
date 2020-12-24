package handlers

import (
	"mathbattle/application"
	"mathbattle/infrastructure"
	"mathbattle/models/mathbattle"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Stat struct {
	Handler
	Replier     application.Replier
	StatService mathbattle.StatService
}

func (h *Stat) Name() string {
	return h.Handler.Name
}

func (h *Stat) Description() string {
	return h.Handler.Description
}

func (h *Stat) IsShowInHelp(ctx infrastructure.TelegramUserContext) bool {
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *Stat) IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error) {
	return true, "", nil
}

func (h *Stat) IsAdminOnly() bool {
	return true
}

func (h *Stat) Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	stat, err := h.StatService.Stat()
	if err != nil {
		return -1, noResponse(), err
	}

	return -1, OneTextResp(h.Replier.FormatStat(stat)), nil
}
