package handlers

import (
	mreplier "mathbattle/application"
	"mathbattle/infrastructure"
	"mathbattle/models/mathbattle"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Unsubscribe struct {
	Handler
	Replier            mreplier.Replier
	ParticipantService mathbattle.ParticipantService
	RoundService       mathbattle.RoundService
}

func (h *Unsubscribe) Name() string {
	return h.Handler.Name
}

func (h *Unsubscribe) Description() string {
	return h.Handler.Description
}

func (h *Unsubscribe) IsShowInHelp(ctx infrastructure.TelegramUserContext) bool {
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *Unsubscribe) IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error) {
	participant, err := h.ParticipantService.GetByTelegramID(ctx.User.TelegramID)
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, h.Replier.NotSubscribed(), nil
		}

		return false, "", err
	}

	if !participant.IsActive {
		return false, h.Replier.NotSubscribed(), nil
	}

	_, err = h.RoundService.GetRunning()
	if err != nil {
		if err != mathbattle.ErrNotFound {
			return false, "", err
		} else {
			return true, "", nil
		}
	}

	return false, "", nil
}

func (h *Unsubscribe) IsAdminOnly() bool {
	return false
}

func (h *Unsubscribe) Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	participant, err := h.ParticipantService.GetByTelegramID(ctx.User.TelegramID)
	if err != nil && err != mathbattle.ErrNotFound {
		return -1, noResponse(), err
	}

	err = h.ParticipantService.Unsubscribe(participant.ID)
	if err != nil {
		return -1, noResponse(), err
	}

	return -1, OneTextResp(h.Replier.UnsubscribeSuccess()), nil
}
