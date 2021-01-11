package handlers

import (
	"mathbattle/application"
	"mathbattle/infrastructure"
	"mathbattle/models/mathbattle"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Subscribe struct {
	Handler
	Replier            application.Replier
	ParticipantService mathbattle.ParticipantService
}

func (h *Subscribe) Name() string {
	return h.Handler.Name
}

func (h *Subscribe) Description() string {
	return h.Handler.Description
}

func (h *Subscribe) IsShowInHelp(ctx infrastructure.TelegramUserContext) bool {
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *Subscribe) IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error) {
	participant, err := h.ParticipantService.GetByTelegramID(ctx.User.TelegramID)
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return true, "", nil
		}

		return false, "", err
	}

	if participant.IsActive {
		return false, h.Replier.AlreadyRegistered(), nil
	}

	return true, "", nil
}

func (h *Subscribe) IsAdminOnly() bool {
	return false
}

func (h *Subscribe) Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	switch ctx.CurrentStep {
	case 0:
		return h.stepCheckExistance(ctx, m)
	case 1:
		return h.stepAcceptName(ctx, m)
	case 2:
		return h.stepAcceptGradeAndFinish(ctx, m)
	default:
		return -1, noResponse(), nil
	}
}

func (h *Subscribe) stepCheckExistance(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	participant, err := h.ParticipantService.GetByTelegramID(ctx.User.TelegramID)
	if err != nil && err != mathbattle.ErrNotFound {
		return -1, noResponse(), err
	}

	if err == mathbattle.ErrNotFound {
		return 1, OneTextResp(h.Replier.RegisterNameExpect()), nil
	}

	participant.IsActive = true
	err = h.ParticipantService.Update(participant)
	if err != nil {
		return -1, noResponse(), err
	}

	return -1, OneTextResp(h.Replier.RegisterSuccess()), nil
}

func (h *Subscribe) stepAcceptName(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	name, ok := mathbattle.ValidateUserName(m.Text)
	if !ok {
		return 1, OneTextResp(h.Replier.RegisterNameWrong()), nil
	}

	ctx.Variables["name"] = infrastructure.NewContextVariableStr(name)

	return 2, OneTextResp(h.Replier.RegisterGradeExpect()), nil
}

func (h *Subscribe) stepAcceptGradeAndFinish(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	grade, ok := mathbattle.ValidateUserGrade(m.Text)
	if !ok {
		return 2, OneTextResp(h.Replier.RegisterGradeWrong()), nil
	}

	_, err := h.ParticipantService.Store(mathbattle.Participant{
		User:     ctx.User,
		Name:     ctx.Variables["name"].AsString(),
		School:   ctx.Variables["school"].AsString(),
		Grade:    grade,
		IsActive: true,
	})

	if err != nil {
		return -1, noResponse(), err
	}

	return -1, OneTextResp(h.Replier.RegisterSuccess()), nil
}
