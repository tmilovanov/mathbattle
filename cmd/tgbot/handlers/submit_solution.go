package handlers

import (
	"io/ioutil"
	"log"
	mreplier "mathbattle/cmd/tgbot/replier"
	mathbattle "mathbattle/models"
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
		if err == mathbattle.ErrRoundNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (h *SubmitSolution) Handle(ctx mathbattle.TelegramUserContext, m *tb.Message) (int, error) {
	r, err := h.Rounds.GetRunning()
	if err != nil {
		return -1, err
	}

	p, exist, err := h.Participants.GetByTelegramID(strconv.FormatInt(ctx.ChatID, 10))
	if err != nil || !exist {
		return -1, err
	}

	log.Printf("SubmitSolution.Handle(), step: %d", ctx.CurrentStep)

	switch ctx.CurrentStep {
	case 0:
		isSuitable, err := h.IsCommandSuitable(ctx)
		if err != nil {
			return -1, err
		}

		if !isSuitable {
			// TODO: Send help message
			ctx.SendText("Not suitable")
			return -1, nil
		}

		return 1, ctx.SendText(h.Replier.GetReply(mreplier.ReplySSolutionExpectProblem))
	case 1: // Expect problem number
		//log.Printf("Round distribution(all): %v", r.ProblemDistribution)
		//log.Printf("Round distribution: %v", r.ProblemDistribution[p.ID])

		problemNumber, err := strconv.Atoi(m.Text)
		if err != nil {
			return 1, ctx.SendText(h.Replier.GetReply(mreplier.ReplySSolutionWrongProblemNumberFormat))
		}
		problemNumber = problemNumber - 1
		if problemNumber < 0 || problemNumber >= len(r.ProblemDistribution[p.ID]) {
			return 1, ctx.SendText(h.Replier.GetReply(mreplier.ReplySSolutionWrongProblemNumber))
		}

		problemId := r.ProblemDistribution[p.ID][problemNumber]
		ctx.Variables["problem_id"] = problemId
		log.Printf("Problem number: %d, id: %s", problemNumber, problemId)
		return 2, nil
	case 2:
		if m.Photo == nil {
			return 2, ctx.SendText(h.Replier.GetReply(mreplier.ReplyWrongSolutionFormat))
		}

		log.Printf("Width: %d, Height: %d Size: %d", m.Photo.Width, m.Photo.Height, m.Photo.FileSize)
		log.Printf("Album ID: %v", m.AlbumID)

		filerReader, err := ctx.Bot.GetFile(&m.Photo.File)
		if err != nil {
			return -1, err
		}

		content, err := ioutil.ReadAll(filerReader)
		if err != nil {
			return -1, err
		}

		s, err := h.Solutions.FindOrCreate(r.ID, p.ID, ctx.Variables["problem_id"])
		if err != nil {
			return -1, err
		}

		err = h.Solutions.AppendPart(s.ID, mathbattle.Image{
			Extension: ".png",
			Content:   content,
		})
		if err != nil {
			return -1, err
		}
	}

	return -1, nil
}
