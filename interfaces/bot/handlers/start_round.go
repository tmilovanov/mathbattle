package handlers

import (
	"errors"

	"mathbattle/application"
	"mathbattle/infrastructure"
	"mathbattle/models/mathbattle"

	tb "gopkg.in/tucnak/telebot.v2"
)

type StartRound struct {
	Handler
	Replier      application.Replier
	RoundService mathbattle.RoundService
}

func (h *StartRound) Name() string {
	return h.Handler.Name
}

func (h *StartRound) Description() string {
	return h.Handler.Description
}

func (h *StartRound) IsShowInHelp(ctx infrastructure.TelegramUserContext) bool {
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *StartRound) IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error) {
	_, err := h.RoundService.GetRunning()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return true, "", nil
		}
		return false, "", err
	}

	return false, "", nil
}

func (h *StartRound) IsAdminOnly() bool {
	return true
}

func (h *StartRound) Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	switch ctx.CurrentStep {
	case 0:
		return h.stepAskDuration(ctx, m)
	case 1:
		return h.stepConfirmDuration(ctx, m)
	case 2:
		return h.stepStart(ctx, m)
	default:
		return -1, noResponse(), nil
	}
}

func (h *StartRound) stepAskDuration(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	return 1, OneTextResp(h.Replier.StartRoundGetDuration()), nil
}

func (h *StartRound) stepConfirmDuration(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	ctx.Variables["until_date"] = infrastructure.NewContextVariableStr(m.Text)
	untilDate, err := mathbattle.ParseStageEndDate(m.Text)
	if err != nil {
		if err == mathbattle.ErrWrongUserInput {
			return 1, OneTextResp(h.Replier.StartRoundWrongDuration()), nil
		}
		return -1, noResponse(), nil
	}

	return 2, OneWithKb(h.Replier.StartRoundConfirmDuration(untilDate), h.Replier.Yes(), h.Replier.No()), nil
}

func (h *StartRound) stepStart(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	if m.Text != h.Replier.Yes() {
		return -1, OneTextResp(h.Replier.Cancel()), nil
	}

	untilDateStr, exist := ctx.Variables["until_date"]
	if !exist {
		return -1, noResponse(), errors.New("Can't find until_date")
	}

	_, err := h.RoundService.StartNew(mathbattle.StartOrder{StageEnd: untilDateStr.AsString()})
	if err != nil {
		return -1, noResponse(), err
	}

	return -1, OneTextResp(h.Replier.StartRoundSuccess()), nil
}
