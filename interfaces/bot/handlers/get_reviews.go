package handlers

import (
	"fmt"
	mreplier "mathbattle/application"
	"mathbattle/infrastructure"
	"mathbattle/models/mathbattle"

	tb "gopkg.in/tucnak/telebot.v2"
)

type GetReviews struct {
	Handler
	Replier            mreplier.Replier
	ParticipantService mathbattle.ParticipantService
	ReviewService      mathbattle.ReviewService
	RoundService       mathbattle.RoundService
	SolutionService    mathbattle.SolutionService
}

func (h *GetReviews) Name() string {
	return h.Handler.Name
}

func (h *GetReviews) Description() string {
	return h.Handler.Description
}

func (h *GetReviews) IsShowInHelp(ctx infrastructure.TelegramUserContext) bool {
	isSuitable, _, _ := h.IsCommandSuitable(ctx)
	return isSuitable
}

func (h *GetReviews) IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error) {
	_, err := h.ParticipantService.GetByTelegramID(ctx.User.TelegramID)
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, h.Replier.NotParticipant(), nil
		}
		return false, "", err
	}

	round, err := h.RoundService.GetLast()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, "", nil
		}
		return false, "", err
	}

	if mathbattle.GetRoundStage(round) != mathbattle.StageFinished {
		return false, "", err
	}

	return true, "", err
}

func (h *GetReviews) IsAdminOnly() bool {
	return false
}

func (h *GetReviews) Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	lastRound, err := h.RoundService.GetLast()
	if err != nil {
		return -1, noResponse(), err
	}

	participant, err := h.ParticipantService.GetByTelegramID(ctx.User.TelegramID)
	if err != nil {
		return -1, noResponse(), err
	}

	solutions, err := h.SolutionService.Find(mathbattle.FindDescriptor{
		RoundID:       lastRound.ID,
		ParticipantID: participant.ID,
		ProblemID:     "",
	})
	if err != nil {
		return -1, noResponse(), err
	}

	result := []TelegramResponse{}
	for _, solution := range solutions {
		reviews, err := h.ReviewService.FindMany(mathbattle.ReviewFindDescriptor{
			ReviewerID: "",
			SolutionID: solution.ID,
		})
		if err != nil {
			return -1, noResponse(), err
		}
		for i, review := range reviews {
			problemCaption := ""
			for _, desc := range lastRound.ProblemDistribution[participant.ID] {
				if desc.ProblemID == solution.ProblemID {
					problemCaption = desc.Caption
				}
			}

			msg := fmt.Sprintf("Комментарий №%d на задачу %s", i+1, problemCaption)
			msg += "\n"
			msg += review.Content
			result = append(result, NewResp(msg))
		}
	}

	return -1, result, nil
}
