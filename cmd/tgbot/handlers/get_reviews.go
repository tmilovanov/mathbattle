package handlers

import (
	"fmt"
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"

	tb "gopkg.in/tucnak/telebot.v2"
)

type GetReviews struct {
	Handler
	Replier     mreplier.Replier
	Participant mathbattle.ParticipantRepository
	Reviews     mathbattle.ReviewRepository
	Rounds      mathbattle.RoundRepository
	Solutions   mathbattle.SolutionRepository
}

func (h *GetReviews) Name() string {
	return h.Handler.Name
}

func (h *GetReviews) Description() string {
	return h.Handler.Description
}

func (h *GetReviews) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	isSuitable, _, _ := h.IsCommandSuitable(ctx)
	return isSuitable
}

func (h *GetReviews) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, string, error) {
	_, err := h.Participant.GetByTelegramID(ctx.User.ChatID)
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, h.Replier.NotParticipant(), nil
		}
		return false, "", err
	}

	round, err := h.Rounds.GetLast()
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

func (h *GetReviews) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	lastRound, err := h.Rounds.GetLast()
	if err != nil {
		return -1, noResponse(), err
	}

	participant, err := h.Participant.GetByTelegramID(ctx.User.ChatID)
	if err != nil {
		return -1, noResponse(), err
	}

	solutions, err := h.Solutions.FindMany(lastRound.ID, participant.ID, "")
	if err != nil {
		return -1, noResponse(), err
	}

	result := []mathbattle.TelegramResponse{}
	i := 1
	for _, solution := range solutions {
		reviews, err := h.Reviews.FindMany("", solution.ID)
		if err != nil {
			return -1, noResponse(), err
		}
		for _, review := range reviews {
			msg := fmt.Sprintf("Комментарий №%d", i)
			msg += "\n"
			msg += review.Content
			result = append(result, mathbattle.NewResp(msg))
		}
	}

	return -1, result, nil
}
