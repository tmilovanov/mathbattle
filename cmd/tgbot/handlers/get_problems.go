package handlers

import (
	"fmt"
	mathbattle "mathbattle/models"

	tb "gopkg.in/tucnak/telebot.v2"
)

type GetProblems struct {
	Handler

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
	res, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *GetProblems) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, error) {
	isReg, err := mathbattle.IsRegistered(h.Participants, ctx.User.ChatID)
	if err != nil {
		return false, err
	}

	if !isReg {
		return false, nil
	}

	_, err = h.Rounds.GetSolveRunning()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
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

	var problemIDs []string
	participantProblems, areExist := curRound.ProblemDistribution[participant.ID]
	if !areExist {
		// New participant
		problems, err := h.ProblemDistributor.GetForParticipant(participant)
		if err != nil {
			return -1, noResponse(), err
		}
		problemIDs = mathbattle.GetProblemIDs(problems)
		curRound.ProblemDistribution[participant.ID] = problemIDs
		err = h.Rounds.Update(curRound)
		if err != nil {
			return -1, noResponse(), err
		}
	} else {
		// Exisiting participant
		problemIDs = participantProblems
	}

	result := []mathbattle.TelegramResponse{}
	for i, problemID := range problemIDs {
		problem, err := h.Problems.GetByID(problemID)
		if err != nil {
			return -1, noResponse(), err
		}

		msg := mathbattle.NewRespImage(mathbattle.Image{
			Extension: problem.Extension,
			Content:   problem.Content,
		})
		msg.Text = fmt.Sprint(i + 1)

		result = append(result, msg)
	}

	return -1, result, nil
}
