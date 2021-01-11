package handlers

import (
	mreplier "mathbattle/application"
	"mathbattle/infrastructure"
	"mathbattle/models/mathbattle"

	tb "gopkg.in/tucnak/telebot.v2"
)

type SubmitReview struct {
	Handler

	Replier            mreplier.Replier
	ReviewService      mathbattle.ReviewService
	ParticipantService mathbattle.ParticipantRepository
	RoundService       mathbattle.RoundService
}

func (h *SubmitReview) Name() string {
	return h.Handler.Name
}

func (h *SubmitReview) Description() string {
	return h.Handler.Description
}

func (h *SubmitReview) IsShowInHelp(ctx infrastructure.TelegramUserContext) bool {
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *SubmitReview) IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error) {
	participant, err := h.ParticipantService.GetByTelegramID(ctx.User.TelegramID)
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, h.Replier.NotParticipant(), nil
		}
		return false, "", err
	}

	round, err := h.RoundService.GetReviewRunning()
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

func (h *SubmitReview) Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	r, err := h.RoundService.GetReviewRunning()
	if err != nil {
		return -1, noResponse(), err
	}

	p, err := h.ParticipantService.GetByTelegramID(ctx.User.TelegramID)
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

func (h *SubmitReview) stepStart(ctx infrastructure.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []TelegramResponse, error) {

	descriptors, err := h.ReviewService.RevewStageDescriptors(participant.ID)
	if err != nil {
		return -1, noResponse(), err
	}

	captions := h.Replier.ReviewGetSolutionCaptions(descriptors)

	return 1, OneWithKb(h.Replier.ReviewExpectSolutionCaption(), captions...), nil
}

func (h *SubmitReview) stepExpectSolutionCaption(ctx infrastructure.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []TelegramResponse, error) {

	descriptors, err := h.ReviewService.RevewStageDescriptors(participant.ID)
	if err != nil {
		return -1, noResponse(), err
	}

	captions := h.Replier.ReviewGetSolutionCaptions(descriptors)

	descriptor, isOk := h.Replier.ReviewGetDescriptor(m.Text)
	if !isOk {
		return 1, OneWithKb(h.Replier.ReviewWrongSolutionCaption(), captions...), nil
	}

	solutionID, isOk := mathbattle.FindSolutionIDbyDescriptor(descriptor, descriptors)
	if !isOk {
		return 1, OneWithKb(h.Replier.ReviewWrongSolutionCaption(), captions...), nil
	}

	ctx.Variables["solution_id"] = infrastructure.NewContextVariableStr(solutionID)
	reviews, err := h.ReviewService.FindMany(mathbattle.ReviewFindDescriptor{
		ReviewerID: participant.ID,
		SolutionID: solutionID,
	})
	if err != nil {
		return -1, noResponse(), err
	}

	if len(reviews) == 0 {
		return 3, OneTextResp(h.Replier.ReviewExpectContent()), nil
	}

	return 2, OneWithKb(h.Replier.ReviewIsRewriteOld(), h.Replier.Yes(), h.Replier.No()), nil
}

func (h *SubmitReview) stepAlreadySubmitted(ctx infrastructure.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []TelegramResponse, error) {

	if m.Text == h.Replier.Yes() {
		solutionID := ctx.Variables["solution_id"].AsString()
		reviews, err := h.ReviewService.FindMany(mathbattle.ReviewFindDescriptor{
			ReviewerID: participant.ID,
			SolutionID: solutionID,
		})
		if err != nil {
			return -1, noResponse(), err
		}

		if len(reviews) != 1 {
			return -1, noResponse(), err
		}

		if err := h.ReviewService.Delete(reviews[0].ID); err != nil {
			return -1, noResponse(), err
		}

		return 3, OneTextResp(h.Replier.ReviewExpectContent()), nil
	} else {
		return -1, OneTextResp(h.Replier.Cancel()), nil
	}
}

func (h *SubmitReview) stepAcceptReview(ctx infrastructure.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []TelegramResponse, error) {

	solutionID := ctx.Variables["solution_id"].AsString()
	_, err := h.ReviewService.Store(mathbattle.Review{
		ReviewerID: participant.ID,
		SolutionID: solutionID,
		Content:    m.Text,
	})
	if err != nil {
		return -1, noResponse(), err
	}

	return -1, OneTextResp(h.Replier.ReviewUploadSuccess()), nil
}
