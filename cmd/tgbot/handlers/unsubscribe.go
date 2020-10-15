package handlers

import (
	"strconv"

	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Unsubscribe struct {
	Handler
	Replier      mreplier.Replier
	Participants mathbattle.ParticipantRepository
	Rounds       mathbattle.RoundRepository
}

func (h *Unsubscribe) Name() string {
	return h.Handler.Name
}

func (h *Unsubscribe) Description() string {
	return h.Handler.Description
}

func (h *Unsubscribe) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	res, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *Unsubscribe) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, error) {
	isReg, err := mathbattle.IsRegistered(h.Participants, ctx.User.ChatID)
	if err != nil {
		return false, err
	}

	if !isReg {
		return false, nil
	}

	_, err = h.Rounds.GetSolveRunning()
	if err != nil {
		if err != mathbattle.ErrNotFound {
			return false, err
		} else {
			return true, nil
		}
	}

	return false, nil
}

func (h *Unsubscribe) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, mathbattle.TelegramResponse, error) {
	var noRepsonse mathbattle.TelegramResponse

	participant, err := h.Participants.GetByTelegramID(strconv.FormatInt(ctx.User.ChatID, 10))
	if err != nil && err != mathbattle.ErrNotFound {
		return -1, noRepsonse, err
	}

	if err == mathbattle.ErrNotFound {
		return -1, mathbattle.NewResp(h.Replier.NotSubscribed()), nil
	}

	err = h.Participants.Delete(participant.ID)
	if err != nil {
		return -1, noRepsonse, err
	}

	return -1, mathbattle.NewResp(h.Replier.UnsubscribeSuccess()), nil
}
