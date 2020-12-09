package handlers

import (
	"io/ioutil"
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"
	"path/filepath"

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
	res, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *SubmitSolution) IsCommandSuitable(ctx mathbattle.TelegramUserContext) (bool, error) {
	isReg, err := mathbattle.IsRegistered(h.Participants, ctx.User.ChatID)
	if err != nil {
		return false, err
	}

	if !isReg {
		return false, nil
	}

	_, err = h.Rounds.GetSolveRunning()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (h *SubmitSolution) IsAdminOnly() bool {
	return false
}

func (h *SubmitSolution) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, []mathbattle.TelegramResponse, error) {
	r, err := h.Rounds.GetSolveRunning()
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
		return h.stepExpectProblemNumber(ctx, m, r, p)
	case 2:
		return h.stepAlreadySubmitted(ctx, m, r, p)
	case 3:
		return h.stepAcceptSolutionPart(ctx, m, r, p)
	default:
		return -1, noResponse(), nil
	}
}

func (h *SubmitSolution) stepStart(ctx mathbattle.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []mathbattle.TelegramResponse, error) {

	captions := mathbattle.ProblemsCaptions(round, participant)
	return 1, mathbattle.OneWithKb(h.Replier.SolutionExpectProblemCaption(), captions...), nil
}

func (h *SubmitSolution) stepExpectProblemNumber(ctx mathbattle.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []mathbattle.TelegramResponse, error) {

	problemDescriptors := round.ProblemDistribution[participant.ID]
	captions := mathbattle.ProblemsCaptions(round, participant)
	problemNumber, isOk := mathbattle.ValidateCaptions(m.Text, problemDescriptors)
	if !isOk {
		return 1, mathbattle.OneWithKb(h.Replier.SolutionWrongProblemCaption(), captions...), nil
	}

	problemID := round.ProblemDistribution[participant.ID][problemNumber].ProblemID
	ctx.Variables["problem_id"] = mathbattle.NewContextVariableStr(problemID)
	ctx.Variables["total_uploaded"] = mathbattle.NewContextVariableInt(0)
	currentSolution, err := h.Solutions.Find(round.ID, participant.ID, problemID)
	if err != nil && err != mathbattle.ErrNotFound {
		return -1, noResponse(), err
	}

	if err == mathbattle.ErrNotFound || len(currentSolution.Parts) == 0 {
		return 3, mathbattle.OneWithKb(h.Replier.SolutionExpectPart(), h.Replier.SolutionFinishUploading()), nil
	}

	return 2, mathbattle.OneWithKb(h.Replier.SolutionIsRewriteOld(), h.Replier.Yes(), h.Replier.No()), nil
}

func (h *SubmitSolution) stepAlreadySubmitted(ctx mathbattle.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []mathbattle.TelegramResponse, error) {

	if m.Text == h.Replier.Yes() {
		problemID := ctx.Variables["problem_id"].AsString()
		solution, err := h.Solutions.Find(round.ID, participant.ID, problemID)
		if err != nil {
			return -1, noResponse(), err
		}

		if err := h.Solutions.Delete(solution.ID); err != nil {
			return -1, noResponse(), err
		}

		return 3, mathbattle.OneWithKb(h.Replier.SolutionExpectPart(), h.Replier.SolutionFinishUploading()), nil
	} else {
		return -1, mathbattle.OneTextResp(h.Replier.SolutionDeclineRewriteOld()), nil
	}
}

func (h *SubmitSolution) stepAcceptSolutionPart(ctx mathbattle.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []mathbattle.TelegramResponse, error) {

	if m.Text == h.Replier.SolutionFinishUploading() {
		totalUploaded, _ := ctx.Variables["total_uploaded"].AsInt()
		if totalUploaded == 0 {
			return -1, mathbattle.OneTextResp(h.Replier.SolutionEmpty()), nil
		} else {
			return -1, mathbattle.OneTextResp(h.Replier.SolutionUploadSuccess(totalUploaded)), nil
		}
	}

	if m.Photo == nil && m.Document == nil {
		return 3, mathbattle.OneWithKb(h.Replier.SolutionWrongFormat(), h.Replier.SolutionFinishUploading()), nil
	}

	var uploadedFile tb.File
	if m.Photo != nil {
		uploadedFile = m.Photo.File
	} else {
		if m.Document != nil {
			uploadedFile = m.Document.File
		}
	}

	content, err := ioutil.ReadAll(uploadedFile.FileReader)
	if err != nil {
		return -1, noResponse(), err
	}
	extension := filepath.Ext(uploadedFile.FilePath)

	s, err := h.Solutions.FindOrCreate(round.ID, participant.ID, ctx.Variables["problem_id"].AsString())
	if err != nil {
		return -1, noResponse(), err
	}

	err = h.Solutions.AppendPart(s.ID, mathbattle.Image{
		Extension: extension,
		Content:   content,
	})
	if err != nil {
		return -1, noResponse(), err
	}

	totalUploaded, err := ctx.Variables["total_uploaded"].AsInt()
	if err != nil {
		return -1, noResponse(), err
	}
	totalUploaded++
	ctx.Variables["total_uploaded"] = mathbattle.NewContextVariableInt(totalUploaded)

	return 3, mathbattle.OneWithKb(h.Replier.SolutionPartUploaded(totalUploaded),
		h.Replier.SolutionFinishUploading()), nil
}
