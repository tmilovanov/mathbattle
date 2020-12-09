package handlers

import (
	"errors"
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"
	"mathbattle/mstd"
	"mathbattle/usecases"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type StartReviewStage struct {
	Handler

	Replier             mreplier.Replier
	Rounds              mathbattle.RoundRepository
	Solutions           mathbattle.SolutionRepository
	Participants        mathbattle.ParticipantRepository
	SolutionDistributor mathbattle.SolutionDistributor
	ReviewersCount      uint
	Postman             mathbattle.TelegramPostman
}

func (h *StartReviewStage) Name() string {
	return h.Handler.Name
}

func (h *StartReviewStage) Description() string {
	return h.Handler.Description
}

func (h *StartReviewStage) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	res, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *StartReviewStage) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, error) {
	_, err := h.Rounds.GetReviewPending()
	if err == nil {
		return true, nil
	}

	return false, err
}

func (h *StartReviewStage) IsAdminOnly() bool {
	return true
}

func (h *StartReviewStage) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	switch ctx.CurrentStep {
	case 0:
		return h.stepConfirmDistribution(ctx, m)
	case 1:
		return h.stepAskDuration(ctx, m)
	case 2:
		return h.stepConfirmDuration(ctx, m)
	case 3:
		return h.stepDistribute(ctx, m)
	default:
		return -1, noResponse(), nil
	}
}

func (h *StartReviewStage) stepConfirmDistribution(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	round, err := h.Rounds.GetReviewPending()
	if err != nil {
		return -1, noResponse(), err
	}

	allRoundSolutions, err := h.Solutions.FindMany(round.ID, "", "")
	if err != nil {
		return -1, noResponse(), err
	}

	distribution := h.SolutionDistributor.Get(allRoundSolutions, h.ReviewersCount)
	distributionDesc, err := usecases.ReviewDistrubitonToString(h.Participants, h.Solutions, distribution)
	if err != nil {
		return -1, noResponse(), err
	}

	return 1, mathbattle.OneWithKb(distributionDesc, h.Replier.Yes(), h.Replier.No()), nil
}

func (h *StartReviewStage) stepAskDuration(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	if m.Text != h.Replier.Yes() {
		return -1, mathbattle.OneTextResp(h.Replier.Cancel()), nil
	}

	return 2, mathbattle.OneTextResp(h.Replier.StartReviewGetDuration()), nil
}

func (h *StartReviewStage) stepConfirmDuration(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	ctx.Variables["until_date"] = mathbattle.NewContextVariableStr(m.Text)
	untilDate, err := mathbattle.ParseStageEndDate(m.Text)
	if err != nil {
		if err == mathbattle.ErrWrongUserInput {
			return 2, mathbattle.OneTextResp(h.Replier.StartReviewWrongDuration()), nil
		}
		return -1, noResponse(), nil
	}

	return 3, mathbattle.OneWithKb(h.Replier.StartReviewConfirmDuration(untilDate), h.Replier.Yes(), h.Replier.No()), nil
}

func (h *StartReviewStage) stepDistribute(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	if m.Text != h.Replier.Yes() {
		return -1, mathbattle.OneTextResp(h.Replier.Cancel()), nil
	}

	untilDateStr, exist := ctx.Variables["until_date"]
	if !exist {
		return -1, noResponse(), errors.New("Can't find until_date")
	}
	untilDate, err := mathbattle.ParseStageEndDate(untilDateStr.AsString())
	if err != nil {
		if err == mathbattle.ErrWrongUserInput {
			return 2, mathbattle.OneTextResp(h.Replier.StartReviewWrongDuration()), nil
		}
		return -1, noResponse(), nil
	}

	round, err := h.Rounds.GetReviewPending()
	if err != nil {
		return -1, noResponse(), err
	}

	allRoundSolutions, err := h.Solutions.FindMany(round.ID, "", "")
	if err != nil {
		return -1, noResponse(), err
	}

	distribution := h.SolutionDistributor.Get(allRoundSolutions, h.ReviewersCount)

	for participantID, solutionIDs := range distribution.BetweenParticipants {
		p, err := h.Participants.GetByID(participantID)
		if err != nil {
			return -1, noResponse(), nil
		}

		err = h.Postman.PostText(p.TelegramID, h.Replier.ReviewPostBefore())
		if err != nil {
			return -1, noResponse(), err
		}

		for i, solutionID := range solutionIDs {
			solution, err := h.Solutions.Get(solutionID)
			if err != nil {
				return -1, noResponse(), err
			}

			problemIndex := mstd.IndexOf(round.ProblemDistribution[p.ID], solution.ProblemID)
			caption := h.Replier.ReviewPostCaption(problemIndex+1, i+1)
			err = h.Postman.PostAlbum(p.TelegramID, caption, solution.Parts)
			if err != nil {
				return -1, noResponse(), err
			}
		}

		err = h.Postman.PostText(p.TelegramID, h.Replier.ReviewPostAfter())
		if err != nil {
			return -1, noResponse(), err
		}
	}

	round.SetReviewStartDate(time.Now())
	round.SetReviewEndDate(untilDate)
	round.ReviewDistribution = distribution
	if err = h.Rounds.Update(round); err != nil {
		return -1, noResponse(), err
	}

	return -1, mathbattle.OneTextResp(h.Replier.StartReviewSuccess()), nil
}
