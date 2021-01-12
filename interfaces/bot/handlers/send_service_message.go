package handlers

import (
	mreplier "mathbattle/application"
	"mathbattle/infrastructure"

	tb "gopkg.in/tucnak/telebot.v2"
)

type SendServiceMessage struct {
	Handler
	Replier mreplier.Replier
}

func (h *SendServiceMessage) Name() string {
	return h.Handler.Name
}

func (h *SendServiceMessage) Description() string {
	return h.Handler.Description
}

func (h *SendServiceMessage) IsShowInHelp(ctx infrastructure.TelegramUserContext) bool {
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *SendServiceMessage) IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error) {
	return true, "", nil
}

func (h *SendServiceMessage) IsAdminOnly() bool {
	return true
}

func (h *SendServiceMessage) Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	switch ctx.CurrentStep {
	case 0:
		return 1, OneTextResp(h.Replier.ServiceMsgGetText()), nil
	case 1:
		return h.stepAcceptTextAndSend(ctx, m)
	default:
		return -1, noResponse(), nil
	}
}

func (h *SendServiceMessage) stepAcceptTextAndSend(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	return -1, noResponse(), nil
}
