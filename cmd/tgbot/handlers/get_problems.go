package handlers

import (
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"
	"mathbattle/mstd"

	tb "gopkg.in/tucnak/telebot.v2"
)

type GetProblems struct {
	Handler

	Replier            mreplier.Replier
	Participants       mathbattle.ParticipantRepository
	Rounds             mathbattle.RoundRepository
	Problems           mathbattle.ProblemRepository
	ProblemDistributor mathbattle.ProblemDistributor
}

func (h *GetProblems) Name() string {
	return h.Handler.Name
}

func (h *GetProblems) Description() string {
	return h.Handler.Description
}

func (h *GetProblems) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *GetProblems) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, string, error) {
	isReg, err := mathbattle.IsRegistered(h.Participants, ctx.User.ChatID)
	if err != nil {
		return false, "", err
	}

	if !isReg {
		return false, h.Replier.NotParticipant(), nil
	}

	_, err = h.Rounds.GetSolveRunning()
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

func (h *GetProblems) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	participant, err := h.Participants.GetByTelegramID(ctx.User.ChatID)
	if err != nil {
		return -1, noResponse(), err
	}

	curRound, err := h.Rounds.GetSolveRunning()
	if err != nil {
		return -1, noResponse(), err
	}

	problemDescriptors, areExist := curRound.ProblemDistribution[participant.ID]
	if !areExist {
		// New participant
		problems, err := h.ProblemDistributor.GetForParticipant(participant)
		if err != nil {
			return -1, noResponse(), err
		}
		for i, problem := range problems {
			curRound.ProblemDistribution[participant.ID] = append(curRound.ProblemDistribution[participant.ID],
				mathbattle.ProblemDescriptor{
					Caption:   mstd.IndexToLetter(i),
					ProblemID: problem.ID,
				})
		}

		err = h.Rounds.Update(curRound)
		if err != nil {
			return -1, noResponse(), err
		}

		problemDescriptors = curRound.ProblemDistribution[participant.ID]
	}

	result := []mathbattle.TelegramResponse{}
	for _, problemDescriptor := range problemDescriptors {
		problem, err := h.Problems.GetByID(problemDescriptor.ProblemID)
		if err != nil {
			return -1, noResponse(), err
		}

		msg := mathbattle.NewRespImage(mathbattle.Image{
			Extension: problem.Extension,
			Content:   problem.Content,
		})
		msg.Text = problemDescriptor.Caption

		result = append(result, msg)
	}

	return -1, result, nil
}
