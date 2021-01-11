package handlers

import (
	"mathbattle/application"
	"mathbattle/infrastructure"
	"mathbattle/models/mathbattle"

	tb "gopkg.in/tucnak/telebot.v2"
)

type GetProblems struct {
	Handler

	Replier            application.Replier
	ParticipantService mathbattle.ParticipantService
	RoundService       mathbattle.RoundService
	ProblemService     mathbattle.ProblemService
}

func (h *GetProblems) Name() string {
	return h.Handler.Name
}

func (h *GetProblems) Description() string {
	return h.Handler.Description
}

func (h *GetProblems) IsShowInHelp(ctx infrastructure.TelegramUserContext) bool {
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *GetProblems) IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error) {
	_, err := h.ParticipantService.GetByTelegramID(ctx.User.TelegramID)
	if err != nil {
		return false, "", err
	}

	_, err = h.RoundService.GetRunning()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, h.Replier.NoRoundRunning(), nil
		}
		return false, "", err
	}

	return true, "", nil
}

func (h *GetProblems) IsAdminOnly() bool {
	return false
}

func (h *GetProblems) Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	participant, err := h.ParticipantService.GetByTelegramID(ctx.User.TelegramID)
	if err != nil {
		return -1, noResponse(), err
	}

	problemDescriptors, err := h.RoundService.GetProblemDescriptors(participant.ID)
	if err != nil {
		return -1, noResponse(), err
	}

	result := []TelegramResponse{}
	for _, problemDescriptor := range problemDescriptors {
		problem, err := h.ProblemService.GetByID(problemDescriptor.ProblemID)
		if err != nil {
			return -1, noResponse(), err
		}

		msg := NewRespImage(mathbattle.Image{
			Extension: problem.Extension,
			Content:   problem.Content,
		})
		msg.Text = problemDescriptor.Caption

		result = append(result, msg)
	}

	return -1, result, nil
}
