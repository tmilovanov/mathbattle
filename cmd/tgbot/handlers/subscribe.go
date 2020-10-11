package handlers

import (
	"strconv"
	"time"

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

func resp(input string) mathbattle.TelegramResponse {
	return mathbattle.TelegramResponse(input)
}

func (h *Subscribe) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, mathbattle.TelegramResponse, error) {
	switch ctx.CurrentStep {
	case 0:
		isReg, err := mathbattle.IsRegistered(h.Participants, ctx.ChatID)
		if err != nil {
			return -1, resp(""), err
		}

		if isReg {
			return -1, resp(h.Replier.GetReply(mreplier.ReplyAlreadyRegistered)), nil
		}

		return 1, resp(h.Replier.GetReply(mreplier.ReplyRegisterNameExpect)), nil
	case 1: //expectName
		name, ok := mathbattle.ValidateUserName(m.Text)
		if !ok {
			return 1, resp(h.Replier.GetReply(mreplier.ReplyRegisterNameWrong)), nil
		}

		ctx.Variables["name"] = mathbattle.NewContextVariableStr(name)

		return 2, resp(h.Replier.GetReply(mreplier.ReplyRegisterGradeExpect)), nil
	case 2: //expectGrade
		grade, ok := mathbattle.ValidateUserGrade(m.Text)
		if !ok {
			return 2, resp(h.Replier.GetReply(mreplier.ReplyRegisterGradeWrong)), nil
		}

		_, err := h.Participants.Store(mathbattle.Participant{
			TelegramID:       strconv.FormatInt(ctx.ChatID, 10),
			Name:             ctx.Variables["name"].AsString(),
			School:           ctx.Variables["school"].AsString(),
			Grade:            grade,
			RegistrationTime: time.Now(),
		})
		if err != nil {
			return -1, resp(""), err
		}

		return -1, resp(h.Replier.GetReply(mreplier.ReplyRegisterSuccess)), nil
	}

	return -1, resp(""), nil
}
