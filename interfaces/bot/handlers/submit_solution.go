package handlers

import (
	"io/ioutil"
	"path/filepath"

	"mathbattle/application"
	"mathbattle/infrastructure"
	"mathbattle/models/mathbattle"

	tb "gopkg.in/tucnak/telebot.v2"
)

type SubmitSolution struct {
	Handler
	Replier            application.Replier
	ParticipantService mathbattle.ParticipantService
	RoundService       mathbattle.RoundService
	SolutionService    mathbattle.SolutionService
}

func (h *SubmitSolution) Name() string {
	return h.Handler.Name
}

func (h *SubmitSolution) Description() string {
	return h.Handler.Description
}

func (h *SubmitSolution) IsShowInHelp(ctx infrastructure.TelegramUserContext) bool {
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *SubmitSolution) IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error) {
	_, err := h.ParticipantService.GetByTelegramID(ctx.User.TelegramID)
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, h.Replier.NotParticipant(), nil
		}
		return false, "", err
	}

	round, err := h.RoundService.GetRunning()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return false, h.Replier.NoRoundRunning(), nil
		}
		return false, "", err
	}

	if mathbattle.GetRoundStage(round) != mathbattle.StageSolve {
		return false, h.Replier.NoRoundRunning(), nil
	}

	return true, "", nil
}

func (h *SubmitSolution) IsAdminOnly() bool {
	return false
}

func (h *SubmitSolution) Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	r, err := h.RoundService.GetRunning()
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
		return h.stepExpectProblemNumber(ctx, m, r, p)
	case 2:
		return h.stepAlreadySubmitted(ctx, m, r, p)
	case 3:
		return h.stepAcceptSolutionPart(ctx, m, r, p)
	default:
		return -1, noResponse(), nil
	}
}

func (h *SubmitSolution) stepStart(ctx infrastructure.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []TelegramResponse, error) {

	descriptors, err := h.SolutionService.GetProblemDescriptors(participant.ID)
	if err != nil {
		return -1, noResponse(), err
	}

	captions := mathbattle.ProblemsCaptions(descriptors)

	return 1, OneWithKb(h.Replier.SolutionExpectProblemCaption(), captions...), nil
}

func (h *SubmitSolution) stepExpectProblemNumber(ctx infrastructure.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []TelegramResponse, error) {

	descriptors, err := h.SolutionService.GetProblemDescriptors(participant.ID)
	if err != nil {
		return -1, noResponse(), err
	}
	captions := mathbattle.ProblemsCaptions(descriptors)

	problemNumber, isOk := mathbattle.ValidateCaptions(m.Text, descriptors)
	if !isOk {
		return 1, OneWithKb(h.Replier.SolutionWrongProblemCaption(), captions...), nil
	}

	problemID := round.ProblemDistribution[participant.ID][problemNumber].ProblemID
	ctx.Variables["problem_id"] = infrastructure.NewContextVariableStr(problemID)
	ctx.Variables["total_uploaded"] = infrastructure.NewContextVariableInt(0)

	solutions, err := h.SolutionService.Find(mathbattle.FindDescriptor{
		RoundID:       round.ID,
		ParticipantID: participant.ID,
		ProblemID:     problemID,
	})
	if err != nil {
		return -1, noResponse(), err
	}

	if len(solutions) == 0 || len(solutions[0].Parts) == 0 {
		_, err = h.SolutionService.Create(mathbattle.Solution{
			ParticipantID: participant.ID,
			ProblemID:     ctx.Variables["problem_id"].AsString(),
			RoundID:       round.ID,
			Mark:          -1,
		})
		if err != nil {
			return -1, noResponse(), err
		}
		return 3, OneWithKb(h.Replier.SolutionExpectPart(), h.Replier.SolutionFinishUploading()), nil
	}

	return 2, OneWithKb(h.Replier.SolutionIsRewriteOld(), h.Replier.Yes(), h.Replier.No()), nil
}

func (h *SubmitSolution) stepAlreadySubmitted(ctx infrastructure.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []TelegramResponse, error) {

	if m.Text == h.Replier.Yes() {
		problemID := ctx.Variables["problem_id"].AsString()
		solutions, err := h.SolutionService.Find(mathbattle.FindDescriptor{
			RoundID:       round.ID,
			ParticipantID: participant.ID,
			ProblemID:     problemID,
		})
		if err != nil {
			return -1, noResponse(), err
		}

		if err := h.SolutionService.Delete(solutions[0].ID); err != nil {
			return -1, noResponse(), err
		}

		_, err = h.SolutionService.Create(mathbattle.Solution{
			ParticipantID: participant.ID,
			ProblemID:     ctx.Variables["problem_id"].AsString(),
			RoundID:       round.ID,
			Mark:          -1,
		})
		if err != nil {
			return -1, noResponse(), err
		}

		return 3, OneWithKb(h.Replier.SolutionExpectPart(), h.Replier.SolutionFinishUploading()), nil
	} else {
		return -1, OneWithKb(h.Replier.SolutionDeclineRewriteOld()), nil
	}
}

func (h *SubmitSolution) stepAcceptSolutionPart(ctx infrastructure.TelegramUserContext, m *tb.Message,
	round mathbattle.Round, participant mathbattle.Participant) (int, []TelegramResponse, error) {

	if m.Text == h.Replier.SolutionFinishUploading() {
		totalUploaded, _ := ctx.Variables["total_uploaded"].AsInt()
		if totalUploaded == 0 {
			// Удалить пустое решение
			problemID := ctx.Variables["problem_id"].AsString()
			s, err := h.SolutionService.Find(mathbattle.FindDescriptor{
				RoundID:       round.ID,
				ParticipantID: participant.ID,
				ProblemID:     problemID,
			})
			if err != nil {
				return -1, noResponse(), err
			}

			err = h.SolutionService.Delete(s[0].ID)
			if err != nil {
				return -1, noResponse(), err
			}

			return -1, OneTextResp(h.Replier.SolutionEmpty()), nil
		} else {
			return -1, OneTextResp(h.Replier.SolutionUploadSuccess(totalUploaded)), nil
		}
	}

	if m.Photo == nil && m.Document == nil {
		return 3, OneWithKb(h.Replier.SolutionWrongFormat(), h.Replier.SolutionFinishUploading()), nil
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

	problemID := ctx.Variables["problem_id"].AsString()
	s, err := h.SolutionService.Find(mathbattle.FindDescriptor{
		RoundID:       round.ID,
		ParticipantID: participant.ID,
		ProblemID:     problemID,
	})
	if err != nil {
		return -1, noResponse(), err
	}
	curSolution := s[0]

	err = h.SolutionService.AppendPart(curSolution.ID, mathbattle.Image{
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
	ctx.Variables["total_uploaded"] = infrastructure.NewContextVariableInt(totalUploaded)

	return 3, OneWithKb(h.Replier.SolutionPartUploaded(totalUploaded),
		h.Replier.SolutionFinishUploading()), nil
}
