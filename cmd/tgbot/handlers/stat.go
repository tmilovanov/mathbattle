package handlers

import (
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"
	"mathbattle/usecases"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Stat struct {
	Handler
	Replier      mreplier.Replier
	Participants mathbattle.ParticipantRepository
	Rounds       mathbattle.RoundRepository
	Solutions    mathbattle.SolutionRepository
	Reviews      mathbattle.ReviewRepository
}

func (h *Stat) Name() string {
	return h.Handler.Name
}

func (h *Stat) Description() string {
	return h.Handler.Description
}

func (h *Stat) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *Stat) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, string, error) {
	return true, "", nil
}

func (h *Stat) IsAdminOnly() bool {
	return true
}

func (h *Stat) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	stat, err := usecases.StatReport(h.Participants, h.Rounds, h.Solutions, h.Reviews)
	if err != nil {
		return -1, noResponse(), err
	}
	return -1, mathbattle.OneTextResp(h.Replier.FormatStat(stat)), nil
}
