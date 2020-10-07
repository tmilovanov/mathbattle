package handlers

import (
	"strconv"

	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Subscribe struct {
	Handler
	Replier      mreplier.Replier
	Participants mathbattle.ParticipantRepository
}

func (h *Subscribe) Name() string {
	return h.Handler.Name
}

func (h *Subscribe) Description() string {
	return h.Handler.Description
}

func (h *Subscribe) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	isReg, err := mathbattle.IsRegistered(h.Participants, ctx.ChatID)
	if err != nil {
		return false
	}

	return !isReg
}

func (h *Subscribe) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, error) {
	switch ctx.CurrentStep {
	case 0:
		isReg, err := mathbattle.IsRegistered(h.Participants, ctx.ChatID)
		if err != nil {
			return -1, err
		}

		if isReg {
			err := ctx.SendText(h.Replier.GetReply(mreplier.ReplyAlreadyRegistered))
			if err != nil {
				return -1, err
			}

			return -1, nil
		}

		err = ctx.SendText(h.Replier.GetReply(mreplier.ReplyRegisterNameExpect))
		if err != nil {
			return -1, err
		}

		return 1, nil
	case 1: //expectName
		name, ok := mathbattle.ValidateUserName(m.Text)
		if !ok {
			err := ctx.SendText(h.Replier.GetReply(mreplier.ReplyRegisterNameWrong))
			if err != nil {
				return -1, err
			}

			return 1, nil
		}

		err := ctx.SendText(h.Replier.GetReply(mreplier.ReplyRegisterGradeExpect))
		if err != nil {
			return -1, err
		}

		ctx.Variables["name"] = mathbattle.NewContextVariableStr(name)
		return 2, nil
	case 2: //expectGrade
		grade, ok := mathbattle.ValidateUserGrade(m.Text)
		if !ok {
			err := ctx.SendText(h.Replier.GetReply(mreplier.ReplyRegisterGradeWrong))
			if err != nil {
				return -1, err
			}

			return 2, nil
		}

		err := h.Participants.Store(mathbattle.Participant{
			TelegramID:       strconv.FormatInt(ctx.ChatID, 10),
			Name:             ctx.Variables["name"].AsString(),
			School:           ctx.Variables["school"].AsString(),
			Grade:            grade,
			RegistrationTime: m.Time(),
		})
		if err != nil {
			return -1, err
		}

		err = ctx.SendText(h.Replier.GetReply(mreplier.ReplyRegisterSuccess))
		if err != nil {
			return -1, err
		}

		return -1, nil
	}

	return -1, nil
}
