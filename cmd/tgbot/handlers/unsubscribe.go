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
}

func (h *Unsubscribe) Name() string {
	return h.Handler.Name
}

func (h *Unsubscribe) Description() string {
	return h.Handler.Description
}

func (h *Unsubscribe) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	isReg, err := mathbattle.IsRegistered(h.Participants, ctx.ChatID)
	if err != nil {
		return false
	}

	return isReg
}

func (h *Unsubscribe) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, error) {
	participant, exist, err := h.Participants.GetByTelegramID(strconv.FormatInt(ctx.ChatID, 10))
	if err != nil {
		return -1, err
	}

	if !exist {
		return -1, ctx.SendText(h.Replier.GetReply(mreplier.ReplyUnsubscribeNotSubscribed))
	}

	err = h.Participants.Delete(participant.ID)
	if err != nil {
		return -1, err
	}

	return -1, ctx.SendText(h.Replier.GetReply(mreplier.ReplyUnsubscribeSuccess))
}
