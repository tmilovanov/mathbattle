package handlers

import (
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"

	tb "gopkg.in/tucnak/telebot.v2"
)

type SubmitReview struct {
	Handler

	Replier      mreplier.Replier
	Participants mathbattle.ParticipantRepository
	Rounds       mathbattle.RoundRepository
	Reviews      mathbattle.ReviewRepository
}

func (h *SubmitReview) Name() string {
	return h.Handler.Name
}

func (h *SubmitReview) Description() string {
	return h.Handler.Description
}

func (h *SubmitReview) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	res, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *SubmitReview) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, error) {
	participant, err := h.Participants.GetByTelegramID(ctx.User.ChatID)
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	round, err := h.Rounds.GetReviewRunning()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	_, isExist := round.ReviewDistribution.BetweenParticipants[participant.ID]
	if !isExist {
		return false, nil
	}

	return true, nil
}

func (h *SubmitReview) IsAdminOnly() bool {
	return false
}

func (h *SubmitReview) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	r, err := h.Rounds.GetReviewRunning()
	if err != nil {
		return -1, noResponse(), err
	}

	p, err := h.Participants.GetByTelegramID(ctx.User.ChatID)
	if err != nil {
		return -1, noResponse(), err
	}

	switch ctx.CurrentStep {
	case 0:
		return h.stepStart(ctx, m, r, p)
	case 1:
		return h.stepExpectSolutionNumber(ctx, m, r, p)
	case 2:
		return h.stepAlreadySubmitted(ctx, m, r, p)
	case 3:
		return h.stepAcceptReview(ctx, m, r, p)
	default:
		return -1, noResponse(), nil
	}
}

func (h *SubmitReview) stepStart(ctx mathbattle.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []mathbattle.TelegramResponse, error) {

	solutionNumbers := mathbattle.SolutionNumbers(round, participant)
	return 1, mathbattle.OneWithKb(h.Replier.ReviewExpectSolutionNumber(), solutionNumbers...), nil
}

func (h *SubmitReview) stepExpectSolutionNumber(ctx mathbattle.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []mathbattle.TelegramResponse, error) {

	solutionIDs := round.ReviewDistribution.BetweenParticipants[participant.ID]
	solutionNumbers := mathbattle.SolutionNumbers(round, participant)
	solutionNumber, isOk := mathbattle.ValidateIndex(m.Text, solutionIDs)
	if !isOk {
		return 1, mathbattle.OneWithKb(h.Replier.ReviewWrongSolutionNumber(), solutionNumbers...), nil
	}

	solutionID := round.ReviewDistribution.BetweenParticipants[participant.ID][solutionNumber]
	ctx.Variables["solution_id"] = mathbattle.NewContextVariableStr(solutionID)
	reviews, err := h.Reviews.FindMany(participant.ID, solutionID)
	if err != nil {
		return -1, noResponse(), err
	}

	if len(reviews) == 0 {
		return 3, mathbattle.OneWithKb(h.Replier.ReviewExpectContent()), nil
	}

	return 2, mathbattle.OneWithKb(h.Replier.ReviewIsRewriteOld(), h.Replier.Yes(), h.Replier.No()), nil
}

func (h *SubmitReview) stepAlreadySubmitted(ctx mathbattle.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []mathbattle.TelegramResponse, error) {

	if m.Text == h.Replier.Yes() {
		solutionID := ctx.Variables["solution_id"].AsString()
		reviews, err := h.Reviews.FindMany(participant.ID, solutionID)
		if err != nil {
			return -1, noResponse(), err
		}

		if len(reviews) != 1 {
			return -1, noResponse(), err
		}

		if err := h.Reviews.Delete(reviews[0].ID); err != nil {
			return -1, noResponse(), err
		}

		return 3, mathbattle.OneWithKb(h.Replier.ReviewExpectContent()), nil
	} else {
		return -1, mathbattle.OneTextResp(h.Replier.Cancel()), nil
	}
}

func (h *SubmitReview) stepAcceptReview(ctx mathbattle.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []mathbattle.TelegramResponse, error) {

	solutionID := ctx.Variables["solution_id"].AsString()
	_, err := h.Reviews.Store(mathbattle.Review{
		ReviewerID: participant.ID,
		SolutionID: solutionID,
		Content:    m.Text,
	})
	if err != nil {
		return -1, noResponse(), err
	}

	return -1, mathbattle.OneTextResp(h.Replier.ReviewUploadSuccess()), nil
}
