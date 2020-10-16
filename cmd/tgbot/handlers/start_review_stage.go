package handlers

import (
	"bytes"
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"

	tb "gopkg.in/tucnak/telebot.v2"
)

type StartReviewStage struct {
	Handler

	Replier             mreplier.Replier
	Rounds              mathbattle.RoundRepository
	Solutions           mathbattle.SolutionRepository
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

func (h *StartReviewStage) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, mathbattle.TelegramResponse, error) {
	switch ctx.CurrentStep {
	case 0:
		return h.stepSendDistribution(ctx, m)
	case 1:
		return h.stepConfirmDistribution(ctx, m)
	case 2:
		return h.stepDistribute(ctx, m)
	default:
		return -1, noResponse(), nil
	}
}

func (h *StartReviewStage) stepSendDistribution(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, mathbattle.TelegramResponse, error) {
	round, err := h.Rounds.GetReviewPending()
	if err != nil {
		return -1, noResponse(), err
	}

	allRoundSolutions, err := h.Solutions.FindMany(round.ID, "", "")
	if err != nil {
		return -1, noResponse(), err
	}

	distribution := h.SolutionDistributor.Get(allRoundSolutions, h.ReviewersCount)

	return 1, mathbattle.NewRespWithKeyboard(distribution.ToString(), h.Replier.Yes(), h.Replier.No()), nil
}

func (h *StartReviewStage) stepConfirmDistribution(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, mathbattle.TelegramResponse, error) {
	if m.Text != h.Replier.Yes() {
		return -1, mathbattle.NewResp(h.Replier.Cancel()), nil
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
	for solutionID, participantIDs := range distribution.BetweenParticipants {
		solution, err := h.Solutions.Get(solutionID)
		if err != nil {
			return -1, noResponse(), nil
		}

		for _, participantID := range participantIDs {
			h.Postman.Post(participantID, &tb.Message{Text: h.Replier.ReviewPost()})
			for _, part := range solution.Parts {
				msg := &tb.Message{
					Photo: &tb.Photo{File: tb.FromReader(bytes.NewReader(part.Content))},
				}
				h.Postman.Post(participantID, msg)
			}
		}
	}

	return -1, noResponse(), nil
}

func (h *StartReviewStage) stepDistribute(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, mathbattle.TelegramResponse, error) {
	return -1, noResponse(), nil
}
