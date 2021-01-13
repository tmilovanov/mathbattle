package handlers

import (
	mreplier "mathbattle/application"
	"mathbattle/infrastructure"
	"mathbattle/models/mathbattle"

	tb "gopkg.in/tucnak/telebot.v2"
)

type StartJuriCommenting struct {
	Handler
	Replier         mreplier.Replier
	RoundService    mathbattle.RoundService
	SolutionService mathbattle.SolutionService
}

func (h *StartJuriCommenting) Name() string {
	return h.Handler.Name
}

func (h *StartJuriCommenting) Description() string {
	return h.Handler.Description
}

func (h *StartJuriCommenting) IsShowInHelp(ctx infrastructure.TelegramUserContext) bool {
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *StartJuriCommenting) IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error) {
	_, err := h.RoundService.GetRunning()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, h.Replier.NoRoundRunning(), nil
		}

		return false, "", err
	}

	return true, "", nil
}

func (h *StartJuriCommenting) IsAdminOnly() bool {
	return true
}

func (h *StartJuriCommenting) Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	switch ctx.CurrentStep {
	case 0:
		round, err := h.RoundService.GetRunning()
		if err != nil {
			return -1, noResponse(), nil
		}

		_, err = h.SolutionService.Find(mathbattle.FindDescriptor{RoundID: round.ID})
		if err != nil {
			return -1, noResponse(), nil
		}

		return -1, noResponse(), nil
		// send solution, with keyboard: "Comment this solution", "Next solution", "Stop"
	default:
		return -1, noResponse(), nil
	}
}
