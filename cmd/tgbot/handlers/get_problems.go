package handlers

import (
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"

	tb "gopkg.in/tucnak/telebot.v2"
)

type GetProblems struct {
	Handler

	Replier      mreplier.Replier
	Participants mathbattle.ParticipantRepository
	Rounds       mathbattle.RoundRepository
	Problems     mathbattle.ProblemRepository
}

func (h *GetProblems) Name() string {
	return h.Handler.Name
}

func (h *GetProblems) Description() string {
	return h.Handler.Description
}

func (h *GetProblems) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	res, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *GetProblems) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, error) {
	isReg, err := mathbattle.IsRegistered(h.Participants, ctx.User.ChatID)
	if err != nil {
		return false, err
	}

	if !isReg {
		return false, nil
	}

	_, err = h.Rounds.GetSolveRunning()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (h *GetProblems) IsAdminOnly() bool {
	return false
}

func (h *GetProblems) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, mathbattle.TelegramResponse, error) {
	return -1, noResponse(), nil
}
