package handlers

import (
	"io/ioutil"
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"
	"path/filepath"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

type SubmitSolution struct {
	Handler
	Replier      mreplier.Replier
	Participants mathbattle.ParticipantRepository
	Rounds       mathbattle.RoundRepository
	Solutions    mathbattle.SolutionRepository
}

func (h *SubmitSolution) Name() string {
	return h.Handler.Name
}

func (h *SubmitSolution) Description() string {
	return h.Handler.Description
}

func (h *SubmitSolution) IsShowInHelp(ctx mathbattle.TelegramUserContext) bool {
	isReg, err := mathbattle.IsRegistered(h.Participants, ctx.ChatID)
	if err != nil {
		return false
	}

	if !isReg {
		return false
	}

	_, err = h.Rounds.GetRunning()
	if err != nil {
		return false
	}

	return true
}

func (h *SubmitSolution) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, error) {
	isReg, err := mathbattle.IsRegistered(h.Participants, ctx.ChatID)
	if err != nil {
		return false, err
	}

	if !isReg {
		return false, nil
	}

	_, err = h.Rounds.GetRunning()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (h *SubmitSolution) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, mathbattle.TelegramResponse, error) {
	var noResponse mathbattle.TelegramResponse

	r, err := h.Rounds.GetRunning()
	if err != nil {
		return -1, noResponse, err
	}

	p, err := h.Participants.GetByTelegramID(strconv.FormatInt(ctx.ChatID, 10))
	if err != nil {
		return -1, noResponse, err
	}

	problemIDs := r.ProblemDistribution[p.ID]
	problemNumbers := []string{}
	for i := 0; i < len(problemIDs); i++ {
		problemNumbers = append(problemNumbers, strconv.Itoa(i+1))
	}

	switch ctx.CurrentStep {
	case 0:
		isSuitable, err := h.IsCommandSuitable(ctx)
		if err != nil {
			return -1, noResponse, err
		}

		if !isSuitable {
			return -1, noResponse, ErrCommandUnavailable
		}

		return 1, mathbattle.NewRespWithKeyboard(h.Replier.GetReply(mreplier.ReplySSolutionExpectProblem), problemNumbers...), nil
	case 1: // Expect problem number
		problemNumber, isOk := mathbattle.ValidateProblemNumber(m.Text, problemIDs)
		if !isOk {
			return 1, mathbattle.NewRespWithKeyboard(h.Replier.GetReply(mreplier.ReplySSolutionWrongProblemNumber), problemNumbers...), nil
		}

		problemID := r.ProblemDistribution[p.ID][problemNumber]
		ctx.Variables["problem_id"] = mathbattle.NewContextVariableStr(problemID)

		return 2, mathbattle.NewRespWithKeyboard(h.Replier.GetReply(mreplier.ReplySSolutionExpectStartAccept),
			h.Replier.GetReply(mreplier.ReplySSoltuionFinishUploading)), nil
	case 2: // Expect solution photos
		if m.Text == h.Replier.GetReply(mreplier.ReplySSoltuionFinishUploading) {
			totalUploaded, _ := ctx.Variables["total_uploaded"].AsInt()
			return -1, mathbattle.NewResp(h.Replier.GetReplySSolutionUploadSuccess(totalUploaded)), nil
		}

		if m.Photo == nil {
			return 2, mathbattle.NewResp(h.Replier.GetReply(mreplier.ReplyWrongSolutionFormat)), nil
		}

		content, err := ioutil.ReadAll(m.Photo.FileReader)
		if err != nil {
			return -1, noResponse, err
		}

		s, err := h.Solutions.FindOrCreate(r.ID, p.ID, ctx.Variables["problem_id"].AsString())
		if err != nil {
			return -1, noResponse, err
		}

		err = h.Solutions.AppendPart(s.ID, mathbattle.Image{
			Extension: filepath.Ext(m.Photo.FilePath),
			Content:   content,
		})
		if err != nil {
			return -1, noResponse, err
		}

		var totalUploaded int
		_, isExist := ctx.Variables["total_uploaded"]
		if !isExist {
			totalUploaded = 1
		} else {
			totalUploaded, _ = ctx.Variables["total_uploaded"].AsInt()
			totalUploaded++
		}
		ctx.Variables["total_uploaded"] = mathbattle.NewContextVariableInt(totalUploaded)

		return 2, mathbattle.NewRespWithKeyboard(h.Replier.GetReplySSolutionPartUploaded(totalUploaded),
			h.Replier.GetReply(mreplier.ReplySSoltuionFinishUploading)), nil
	}

	return -1, noResponse, nil
}
