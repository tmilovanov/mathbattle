package handlers

import (
	mathbattle "mathbattle/models"

	tb "gopkg.in/tucnak/telebot.v2"
)

type StartReviewStage struct {
	Handler

	Rounds mathbattle.RoundRepository
}

func (h *StartReviewStage) Name() string {
	return h.Handler.Name
}

func (h *StartReviewStage) Description() string {
	return h.Handler.Description
}

func (h *StartReviewStage) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	return false
}

func (h *StartReviewStage) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, error) {
	_, err := h.Rounds.GetReviewPending()
	if err == nil {
		return true, nil
	}

	return false, err
}

func (h *StartReviewStage) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, mathbattle.TelegramResponse, error) {
	switch ctx.CurrentStep {
	case 0:
		return h.stepConfirmDistribution(ctx, m)
	case 1:
		return h.stepDistribute(ctx, m)
	default:
		return -1, noResponse(), nil
	}
}

func (h *StartReviewStage) stepConfirmDistribution(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, mathbattle.TelegramResponse, error) {
	return -1, noResponse(), nil
}

func (h *StartReviewStage) stepDistribute(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, mathbattle.TelegramResponse, error) {
	return -1, noResponse(), nil
}
