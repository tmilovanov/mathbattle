package application

import (
	"errors"
	"fmt"
	"log"
	"time"

	"mathbattle/libs/mstd"
	"mathbattle/models/mathbattle"
)

func ReviewDistrubitonToString(participants mathbattle.ParticipantRepository, solutions mathbattle.SolutionRepository,
	d mathbattle.ReviewDistribution) (mathbattle.ReviewDistributionDesc, error) {

	result := ""
	result += "To orgs: \n"
	result += "----------\n"
	for _, solutionID := range d.ToOrganizers {
		solution, err := solutions.Get(solutionID)
		if err != nil {
			return mathbattle.ReviewDistributionDesc{Desc: ""}, err
		}

		p, err := participants.GetByID(solution.ParticipantID)
		if err != nil {
			return mathbattle.ReviewDistributionDesc{Desc: ""}, err
		}

		result += fmt.Sprintf("'%s', %d grade, solution on %s\n", p.Name, p.Grade, solution.ProblemID)
	}
	result += "\n"

	result += "Between participants: \n"
	result += "---------\n"
	for participantID, solutionIDs := range d.BetweenParticipants {
		p, err := participants.GetByID(participantID)
		if err != nil {
			return mathbattle.ReviewDistributionDesc{Desc: ""}, err
		}

		for _, solutionID := range solutionIDs {
			solution, err := solutions.Get(solutionID)
			if err != nil {
				return mathbattle.ReviewDistributionDesc{Desc: ""}, err
			}
			fromParticipant, err := participants.GetByID(solution.ParticipantID)
			if err != nil {
				return mathbattle.ReviewDistributionDesc{Desc: ""}, err
			}
			result += fmt.Sprintf("Participant '%s' <- Participant '%s' (Problem %s)\n", p.Name, fromParticipant.Name, solution.ProblemID)
		}
	}

	return mathbattle.ReviewDistributionDesc{Desc: result}, nil
}

type ProblemDistributor interface {
	GetForParticipant(participant mathbattle.Participant) ([]mathbattle.Problem, error)
	GetForParticipantCount(participant mathbattle.Participant, count int) ([]mathbattle.Problem, error)
}

// SolutionDistributor распределяет решения участников на ревью после заврешения этапа решения
type SolutionDistributor interface {
	// Распределить все решения, сданные в текущем раунде, на ревью между участниками.
	// Каждое решение будет отправлено нескольким другим участникам (reviewerCount)
	Get(allRoundSolutions []mathbattle.Solution, reviewerCount uint) mathbattle.ReviewDistribution
}

type RoundService struct {
	Rep                    mathbattle.RoundRepository
	Replier                Replier
	Postman                mathbattle.PostmanService
	Participants           mathbattle.ParticipantRepository
	Solutions              mathbattle.SolutionRepository
	Reviews                mathbattle.ReviewRepository
	SolveStageDistributor  ProblemDistributor
	ReviewStageDistributor SolutionDistributor
	ReviewersCount         int
}

func (rs *RoundService) StartNew(startOrder mathbattle.StartOrder) (mathbattle.Round, error) {
	_, err := rs.Rep.GetRunning()
	if err != mathbattle.ErrNotFound {
		if err == nil {
			return mathbattle.Round{}, errors.New("Round already started")
		}
		return mathbattle.Round{}, err
	}

	solveEndTime, err := mathbattle.ParseStageEndDate(startOrder.StageEnd)
	if err != nil {
		log.Printf("Failed to parse stage end date: '%s', Error: '%v'", startOrder.StageEnd, err)
		return mathbattle.Round{}, err
	}

	round := mathbattle.NewRoundFromEnd(solveEndTime)

	participants, err := rs.Participants.GetAll()
	if err != nil {
		return round, err
	}

	for _, participant := range participants {
		participantProblems, err := rs.SolveStageDistributor.GetForParticipant(participant)
		if err != nil {
			return round, err
		}

		for i, problem := range participantProblems {
			round.ProblemDistribution[participant.ID] = append(round.ProblemDistribution[participant.ID],
				mathbattle.ProblemDescriptor{
					Caption:   mstd.IndexToLetter(i),
					ProblemID: problem.ID,
				})
		}

		duration := round.GetSolveStageDuration()
		stageEndMsk, err := round.GetSolveEndDateMsk()
		if err != nil {
			return round, err
		}

		message := rs.Replier.ProblemsPostBefore(duration, stageEndMsk)
		err = rs.Postman.SendSimpleMessage(participant.TelegramID, message)
		if err != nil {
			return round, err
		}

		for i := 0; i < len(participantProblems); i++ {
			err = rs.Postman.SendImage(participant.TelegramID, round.ProblemDistribution[participant.ID][i].Caption,
				participantProblems[i].Content)
			if err != nil {
				return round, err
			}
		}

		err = rs.Postman.SendSimpleMessage(participant.TelegramID, rs.Replier.ProblemsPostAfter())
		if err != nil {
			return round, err
		}
	}

	round, err = rs.Rep.Store(round)
	if err != nil {
		return round, err
	}

	err = rs.StartSchedulingActions()
	if err != nil {
		return round, err
	}

	return round, nil
}

func (rs *RoundService) StartReviewStage(startOrder mathbattle.StartOrder) (mathbattle.Round, error) {
	round := mathbattle.Round{}

	untilDate, err := mathbattle.ParseStageEndDate(startOrder.StageEnd)
	if err != nil {
		return round, err
	}

	round, err = rs.Rep.GetReviewPending()
	if err != nil {
		return round, err
	}

	allRoundSolutions, err := rs.Solutions.FindMany(round.ID, "", "")
	if err != nil {
		return round, err
	}

	distribution := rs.ReviewStageDistributor.Get(allRoundSolutions, uint(rs.ReviewersCount))

	round.SetReviewStartDate(time.Now())
	round.SetReviewEndDate(untilDate)
	round.ReviewDistribution = distribution
	if err = rs.Rep.Update(round); err != nil {
		return round, err
	}

	for participantID, solutionIDs := range distribution.BetweenParticipants {
		p, err := rs.Participants.GetByID(participantID)
		if err != nil {
			return round, err
		}

		endMsk, err := round.GetReviewEndDateMsk()
		if err != nil {
			return round, err
		}

		err = rs.Postman.SendSimpleMessage(p.TelegramID, rs.Replier.ReviewPostBefore(round.GetReviewStageDuration(), endMsk))
		if err != nil {
			return round, err
		}

		descriptors, err := mathbattle.SolutionDescriptorsFromSolutionIDs(rs.Solutions, participantID, round)
		if err != nil {
			return round, err
		}

		for i := 0; i < len(solutionIDs); i++ {
			solution, err := rs.Solutions.Get(solutionIDs[i])
			if err != nil {
				return round, err
			}

			caption := rs.Replier.ReviewPostCaption(descriptors[i].ProblemCaption, descriptors[i].SolutionNumber)

			images := [][]byte{}
			for _, part := range solution.Parts {
				images = append(images, part.Content)
			}

			err = rs.Postman.SendAlbum(p.TelegramID, caption, images)
			if err != nil {
				return round, err
			}
		}

		err = rs.Postman.SendSimpleMessage(p.TelegramID, rs.Replier.ReviewPostAfter())
		if err != nil {
			return round, err
		}
	}

	err = rs.StartSchedulingActions()
	if err != nil {
		return round, err
	}

	return round, nil
}

func (rs *RoundService) ReviewStageDistributionDesc() (mathbattle.ReviewDistributionDesc, error) {
	round, err := rs.Rep.GetReviewPending()
	if err != nil {
		return mathbattle.ReviewDistributionDesc{Desc: ""}, err
	}

	allRoundSolutions, err := rs.Solutions.FindMany(round.ID, "", "")
	if err != nil {
		return mathbattle.ReviewDistributionDesc{Desc: ""}, err
	}

	distribution := rs.ReviewStageDistributor.Get(allRoundSolutions, 2)

	return ReviewDistrubitonToString(rs.Participants, rs.Solutions, distribution)
}

func (rs *RoundService) GetAll() ([]mathbattle.Round, error) {
	return rs.Rep.GetAll()
}

func (rs *RoundService) GetByID(ID string) (mathbattle.Round, error) {
	return rs.Rep.Get(ID)
}

func (rs *RoundService) GetRunning() (mathbattle.Round, error) {
	return rs.Rep.GetRunning()
}

func (rs *RoundService) GetReviewPending() (mathbattle.Round, error) {
	return rs.Rep.GetReviewPending()
}

func (rs *RoundService) GetReviewRunning() (mathbattle.Round, error) {
	return rs.Rep.GetReviewRunning()
}

func (rs *RoundService) GetLast() (mathbattle.Round, error) {
	return rs.Rep.GetLast()
}

func (rs *RoundService) GetProblemDescriptors(participantID string) ([]mathbattle.ProblemDescriptor, error) {
	curRound, err := rs.Rep.GetRunning()
	if err != nil {
		return []mathbattle.ProblemDescriptor{}, err
	}

	participant, err := rs.Participants.GetByID(participantID)
	if err != nil {
		return []mathbattle.ProblemDescriptor{}, err
	}

	problemDescriptors, areExist := curRound.ProblemDistribution[participantID]
	if !areExist {
		// New participant
		problems, err := rs.SolveStageDistributor.GetForParticipant(participant)
		if err != nil {
			return []mathbattle.ProblemDescriptor{}, err
		}

		for i, problem := range problems {
			curRound.ProblemDistribution[participant.ID] = append(curRound.ProblemDistribution[participant.ID],
				mathbattle.ProblemDescriptor{
					Caption:   mstd.IndexToLetter(i),
					ProblemID: problem.ID,
				})
		}

		err = rs.Rep.Update(curRound)
		if err != nil {
			return []mathbattle.ProblemDescriptor{}, err
		}

		problemDescriptors = curRound.ProblemDistribution[participant.ID]
	}

	return problemDescriptors, nil
}

func (rs *RoundService) onSolveStageEnd() {
	participants, err := rs.Participants.GetAll()
	if err != nil {
		log.Printf("onSolveStageEnd - failed to get participants, error: %v", err)
		return
	}

	round, err := rs.Rep.GetRunning()
	if err != nil {
		log.Printf("onSolveStageEnd - failed to get current round, error: %v", err)
		return
	}

	for _, participant := range participants {
		allParticipantSolutions, err := rs.Solutions.FindMany(round.ID, participant.ID, "")
		if err != nil {
			log.Printf("onSolveStageEnd - failed to get all participant solutions, error: %v", err)
		}

		var msg string
		if len(allParticipantSolutions) == 0 {
			msg = rs.Replier.SolveStageEndNoSolutions()
		} else {
			msg = rs.Replier.SolveStageEnd()
		}

		err = rs.Postman.SendSimpleMessage(participant.TelegramID, msg)
		if err != nil {
			log.Printf("onSolveStageEnd - failed to send message to participant: %v", err)
		}
	}
}

func (rs *RoundService) onReviewStageEnd() {
	participants, err := rs.Participants.GetAll()
	if err != nil {
		log.Printf("onReviewStageEnd - failed to get all participants, error: %v", err)
		return
	}

	for _, participant := range participants {
		err := rs.Postman.SendSimpleMessage(participant.TelegramID, rs.Replier.ReviewStageEnd())
		if err != nil {
			log.Printf("onReviewStageEnd - failed to send message to participant, error: %v", err)
		}
	}
}

func (rs *RoundService) StartSchedulingActions() error {
	log.Printf("StartSchedulingActions()")

	round, err := rs.Rep.GetRunning()
	if err != nil {
		if err == mathbattle.ErrNotFound {
			return nil
		}
		log.Printf("StartSchedulingActions(), failed to get current round, error: %v", err)
		return err
	}

	roundStage := mathbattle.GetRoundStage(round)
	log.Printf("StartSchedulingActions(), round stage is %v", roundStage)
	switch roundStage {
	case mathbattle.StageSolve:
		runFuncAfter := time.Until(round.GetSolveEndDate())
		time.AfterFunc(runFuncAfter, rs.onSolveStageEnd)
		log.Printf("StartSchedulingActions(), onSolveStagEnd is scheduled after %v, solve stage end date is %v",
			runFuncAfter, round.GetSolveEndDate())
	case mathbattle.StageReview:
		runFuncAfter := time.Until(round.GetReviewEndDate())
		time.AfterFunc(runFuncAfter, rs.onReviewStageEnd)
		log.Printf("StartSchedulingActions(), onReviewStagEnd is scheduled after %v, solve stage end date is %v",
			runFuncAfter, round.GetReviewEndDate())
	default:
		log.Printf("StartSchedulingActions(), not scheduling anything")
	}

	return nil
}
