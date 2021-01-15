package handlers

import (
	mreplier "mathbattle/application"
	"mathbattle/infrastructure"
	"mathbattle/models/mathbattle"

	tb "gopkg.in/tucnak/telebot.v2"
)

type GetMyResults struct {
	Handler
	Replier            mreplier.Replier
	RoundService       mathbattle.RoundService
	SolutionService    mathbattle.SolutionService
	ParticipantService mathbattle.ParticipantService
	ReviewService      mathbattle.ReviewService
}

func (h *GetMyResults) Name() string {
	return h.Handler.Name
}

func (h *GetMyResults) Description() string {
	return h.Handler.Description
}

func (h *GetMyResults) IsShowInHelp(ctx infrastructure.TelegramUserContext) bool {
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *GetMyResults) IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error) {
	_, err := h.ParticipantService.GetByTelegramID(ctx.User.TelegramID)
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, h.Replier.NotParticipant(), nil
		}
		return false, "", err
	}

	_, err = h.RoundService.GetRunning()
	if err != mathbattle.ErrNotFound {
		return false, "", nil
	}

	_, err = h.RoundService.GetLast()
	if err != nil {
		return false, "", nil
	}

	return true, "", nil
}

func (h *GetMyResults) IsAdminOnly() bool {
	return false
}

func (h *GetMyResults) Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	switch ctx.CurrentStep {
	case 0:
		round, err := h.RoundService.GetLast()
		if err != nil {
			return -1, noResponse(), nil
		}

		participant, err := h.ParticipantService.GetByTelegramID(ctx.User.TelegramID)
		if err != nil {
			return -1, noResponse(), nil
		}

		allResps := []string{}

		for _, problemDesc := range round.ProblemDistribution[participant.ID] {
			solutions, err := h.SolutionService.Find(mathbattle.FindDescriptor{
				RoundID:       round.ID,
				ParticipantID: participant.ID,
				ProblemID:     problemDesc.ProblemID,
			})
			if err != nil {
				return -1, noResponse(), nil
			}

			if len(solutions) == 0 {
				allResps = append(allResps, h.Replier.MyResultsProblemNotSolved(problemDesc.Caption))
			}
			solution := solutions[0]

			otherParticipantReviews, err := h.ReviewService.FindMany(mathbattle.ReviewFindDescriptor{
				SolutionID: solution.ID,
			})
			if err != nil {
				return -1, noResponse(), nil
			}

			response := h.Replier.MyResultsProblemResults(problemDesc.Caption, solution.JuriComment, solution.Mark,
				otherParticipantReviews)

			allResps = append(allResps, response)
		}

		reviewDescriptors, err := h.ReviewService.RevewStageDescriptors(participant.ID)
		if err != nil {
			return -1, noResponse(), err
		}

		for _, desc := range reviewDescriptors {
			reviews, err := h.ReviewService.FindMany(mathbattle.ReviewFindDescriptor{
				ReviewerID: participant.ID,
				SolutionID: desc.SolutionID,
			})
			if err != nil {
				return -1, noResponse(), nil
			}

			if len(reviews) != 0 {
				solution, err := h.SolutionService.Get(desc.SolutionID)
				if err != nil {
					return -1, noResponse(), nil
				}

				allResps = append(allResps,
					h.Replier.MyResultsReviewResults(desc.ProblemCaption, desc.SolutionNumber, true,
						solution.JuriComment, reviews[0].Mark))
			} else {
				allResps = append(allResps,
					h.Replier.MyResultsReviewResults(desc.ProblemCaption, desc.SolutionNumber, false, "", -1))
			}

		}

		return -1, NewResps(allResps...), nil
	default:
		return -1, noResponse(), nil
	}
}
