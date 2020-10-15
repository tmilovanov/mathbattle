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
	res, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *Subscribe) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, error) {
	isReg, err := mathbattle.IsRegistered(h.Participants, ctx.User.ChatID)
	if err != nil {
		return false, err
	}

	return !isReg, nil
}

func (h *Subscribe) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, mathbattle.TelegramResponse, error) {
	var noResponse mathbattle.TelegramResponse

	switch ctx.CurrentStep {
	case 0:
		isReg, err := mathbattle.IsRegistered(h.Participants, ctx.User.ChatID)
		if err != nil {
			return -1, noResponse, err
		}

		if isReg {
			return -1, mathbattle.NewResp(h.Replier.AlreadyRegistered()), nil
		}

		return 1, mathbattle.NewResp(h.Replier.RegisterNameExpect()), nil
	case 1: //expectName
		name, ok := mathbattle.ValidateUserName(m.Text)
		if !ok {
			return 1, mathbattle.NewResp(h.Replier.RegisterNameWrong()), nil
		}

		ctx.Variables["name"] = mathbattle.NewContextVariableStr(name)

		return 2, mathbattle.NewResp(h.Replier.RegisterGradeExpect()), nil
	case 2: //expectGrade
		grade, ok := mathbattle.ValidateUserGrade(m.Text)
		if !ok {
			return 2, mathbattle.NewResp(h.Replier.RegisterGradeWrong()), nil
		}

		_, err := h.Participants.Store(mathbattle.Participant{
			TelegramID:       strconv.FormatInt(ctx.User.ChatID, 10),
			Name:             ctx.Variables["name"].AsString(),
			School:           ctx.Variables["school"].AsString(),
			Grade:            grade,
			RegistrationTime: time.Now(),
		})
		if err != nil {
			return -1, mathbattle.NewResp(""), err
		}

		return -1, mathbattle.NewResp(h.Replier.RegisterSuccess()), nil
	}

	return -1, mathbattle.NewResp(""), nil
}
