package handlers

import (
	"errors"

	"mathbattle/application"
	"mathbattle/infrastructure"
	"mathbattle/models/mathbattle"

	tb "gopkg.in/tucnak/telebot.v2"
)

type StartReviewStage struct {
	Handler

	Replier      application.Replier
	RoundService mathbattle.RoundService
}

func (h *StartReviewStage) Name() string {
	return h.Handler.Name
}

func (h *StartReviewStage) Description() string {
	return h.Handler.Description
}

func (h *StartReviewStage) IsShowInHelp(ctx infrastructure.TelegramUserContext) bool {
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *StartReviewStage) IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error) {
	_, err := h.RoundService.GetReviewPending()
	if err == nil {
		return true, "", nil
	}

	return false, "", err
}

func (h *StartReviewStage) IsAdminOnly() bool {
	return true
}

func (h *StartReviewStage) Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	switch ctx.CurrentStep {
	case 0:
		return h.stepConfirmDistribution(ctx, m)
	case 1:
		return h.stepAskDuration(ctx, m)
	case 2:
		return h.stepConfirmDuration(ctx, m)
	case 3:
		return h.stepDistribute(ctx, m)
	default:
		return -1, noResponse(), nil
	}
}

func (h *StartReviewStage) stepConfirmDistribution(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	distribution, err := h.RoundService.ReviewStageDistributionDesc()
	if err != nil {
		return -1, noResponse(), err
	}

	return 1, OneWithKb(distribution.Desc, h.Replier.Yes(), h.Replier.No()), nil
}

func (h *StartReviewStage) stepAskDuration(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	if m.Text != h.Replier.Yes() {
		return -1, OneTextResp(h.Replier.Cancel()), nil
	}

	return 2, OneTextResp(h.Replier.StartReviewGetDuration()), nil
}

func (h *StartReviewStage) stepConfirmDuration(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	ctx.Variables["until_date"] = infrastructure.NewContextVariableStr(m.Text)
	untilDate, err := mathbattle.ParseStageEndDate(m.Text)
	if err != nil {
		if err == mathbattle.ErrWrongUserInput {
			return 2, OneTextResp(h.Replier.StartReviewWrongDuration()), nil
		}
		return -1, noResponse(), nil
	}

	return 3, OneWithKb(h.Replier.StartReviewConfirmDuration(untilDate), h.Replier.Yes(), h.Replier.No()), nil
}

func (h *StartReviewStage) stepDistribute(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	if m.Text != h.Replier.Yes() {
		return -1, OneTextResp(h.Replier.Cancel()), nil
	}

	untilDateStr, exist := ctx.Variables["until_date"]
	if !exist {
		return -1, noResponse(), errors.New("Can't find until_date")
	}

	_, err := h.RoundService.StartReviewStage(mathbattle.StartOrder{StageEnd: untilDateStr.AsString()})
	if err != nil {
		return -1, noResponse(), err
	}

	return -1, OneTextResp(h.Replier.StartReviewSuccess()), nil
}
