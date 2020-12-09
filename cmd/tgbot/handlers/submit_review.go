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
	Solutions    mathbattle.SolutionRepository
	Postman      mathbattle.TelegramPostman
}

func (h *SubmitReview) Name() string {
	return h.Handler.Name
}

func (h *SubmitReview) Description() string {
	return h.Handler.Description
}

func (h *SubmitReview) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *SubmitReview) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, string, error) {
	participant, err := h.Participants.GetByTelegramID(ctx.User.ChatID)
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, h.Replier.NotParticipant(), nil
		}
		return false, "", err
	}

	round, err := h.Rounds.GetReviewRunning()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, h.Replier.NoRoundRunning(), nil
		}
		return false, "", err
	}

	_, isExist := round.ReviewDistribution.BetweenParticipants[participant.ID]
	if !isExist {
		return false, "", nil
	}

	return true, "", nil
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
		return h.stepExpectSolutionCaption(ctx, m, r, p)
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

	descriptors, err := mathbattle.SolutionDescriptorsFromSolutionIDs(h.Solutions, participant.ID, round)
	if err != nil {
		return -1, noResponse(), err
	}

	captions := h.Replier.ReviewGetSolutionCaptions(descriptors)

	return 1, mathbattle.OneWithKb(h.Replier.ReviewExpectSolutionCaption(), captions...), nil
}

func (h *SubmitReview) stepExpectSolutionCaption(ctx mathbattle.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []mathbattle.TelegramResponse, error) {

	descriptors, err := mathbattle.SolutionDescriptorsFromSolutionIDs(h.Solutions, participant.ID, round)
	if err != nil {
		return -1, noResponse(), err
	}

	captions := h.Replier.ReviewGetSolutionCaptions(descriptors)

	descriptor, isOk := h.Replier.ReviewGetDescriptor(m.Text)
	if !isOk {
		return 1, mathbattle.OneWithKb(h.Replier.ReviewWrongSolutionCaption(), captions...), nil
	}

	solutionID, isOk := mathbattle.FindSolutionIDbyDescriptor(descriptor, descriptors)
	if !isOk {
		return 1, mathbattle.OneWithKb(h.Replier.ReviewWrongSolutionCaption(), captions...), nil
	}

	ctx.Variables["solution_id"] = mathbattle.NewContextVariableStr(solutionID)
	reviews, err := h.Reviews.FindMany(participant.ID, solutionID)
	if err != nil {
		return -1, noResponse(), err
	}

	if len(reviews) == 0 {
		return 3, mathbattle.OneTextResp(h.Replier.ReviewExpectContent()), nil
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

		return 3, mathbattle.OneTextResp(h.Replier.ReviewExpectContent()), nil
	} else {
		return -1, mathbattle.OneTextResp(h.Replier.Cancel()), nil
	}
}

func (h *SubmitReview) stepAcceptReview(ctx mathbattle.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []mathbattle.TelegramResponse, error) {

	solutionID := ctx.Variables["solution_id"].AsString()
	review, err := h.Reviews.Store(mathbattle.Review{
		ReviewerID: participant.ID,
		SolutionID: solutionID,
		Content:    m.Text,
	})
	if err != nil {
		return -1, noResponse(), err
	}

	solution, err := h.Solutions.Get(solutionID)
	if err != nil {
		return -1, noResponse(), err
	}

	reviewedParticipant, err := h.Participants.GetByID(solution.ParticipantID)
	if err != nil {
		return -1, noResponse(), err
	}

	err = h.Postman.PostText(reviewedParticipant.TelegramID, h.Replier.ReviewMsgForReviewee(review))
	if err != nil {
		return -1, noResponse(), err
	}

	return -1, mathbattle.OneTextResp(h.Replier.ReviewUploadSuccess()), nil
}
