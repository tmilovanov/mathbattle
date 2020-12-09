package handlers

import (
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
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *Subscribe) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, string, error) {
	isReg, err := mathbattle.IsRegistered(h.Participants, ctx.User.ChatID)
	if err != nil {
		return false, "", err
	}

	if isReg {
		return false, h.Replier.AlreadyRegistered(), nil
	}

	return true, "", nil
}

func (h *Subscribe) IsAdminOnly() bool {
	return false
}

func (h *Subscribe) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	switch ctx.CurrentStep {
	case 0:
		return 1, mathbattle.OneTextResp(h.Replier.RegisterNameExpect()), nil
	case 1:
		return h.stepAcceptName(ctx, m)
	case 2:
		return h.stepAcceptGradeAndFinish(ctx, m)
	default:
		return -1, noResponse(), nil
	}
}

func (h *Subscribe) stepAcceptName(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	name, ok := mathbattle.ValidateUserName(m.Text)
	if !ok {
		return 1, mathbattle.OneTextResp(h.Replier.RegisterNameWrong()), nil
	}

	ctx.Variables["name"] = mathbattle.NewContextVariableStr(name)

	return 2, mathbattle.OneTextResp(h.Replier.RegisterGradeExpect()), nil
}

func (h *Subscribe) stepAcceptGradeAndFinish(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	grade, ok := mathbattle.ValidateUserGrade(m.Text)
	if !ok {
		return 2, mathbattle.OneTextResp(h.Replier.RegisterGradeWrong()), nil
	}

	_, err := h.Participants.Store(mathbattle.Participant{
		TelegramID:       ctx.User.ChatID,
		Name:             ctx.Variables["name"].AsString(),
		School:           ctx.Variables["school"].AsString(),
		Grade:            grade,
		RegistrationTime: time.Now(),
	})
	if err != nil {
		return -1, noResponse(), err
	}

	return -1, mathbattle.OneTextResp(h.Replier.RegisterSuccess()), nil
}
