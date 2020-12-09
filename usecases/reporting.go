package usecases

import (
	"fmt"
	mathbattle "mathbattle/models"
	"time"
)

func ReviewDistrubitonToString(participants mathbattle.ParticipantRepository, solutions mathbattle.SolutionRepository,
	d mathbattle.ReviewDistribution) (string, error) {

	result := ""
	result += "To orgs: \n"
	result += "----------\n"
	for _, solutionID := range d.ToOrganizers {
		solution, err := solutions.Get(solutionID)
		if err != nil {
			return "", err
		}

		p, err := participants.GetByID(solution.ParticipantID)
		if err != nil {
			return "", err
		}

		result += fmt.Sprintf("'%s', %d grade, solution on %s\n", p.Name, p.Grade, solution.ProblemID)
	}
	result += "\n"

	result += "Between participants: \n"
	result += "---------\n"
	for participantID, solutionIDs := range d.BetweenParticipants {
		p, err := participants.GetByID(participantID)
		if err != nil {
			return "", err
		}

		for _, solutionID := range solutionIDs {
			solution, err := solutions.Get(solutionID)
			if err != nil {
				return "", err
			}
			fromParticipant, err := participants.GetByID(solution.ParticipantID)
			if err != nil {
				return "", err
			}
			result += fmt.Sprintf("Participant '%s' <- Participant '%s' (Problem %s)\n", p.Name, fromParticipant.Name, solution.ProblemID)
		}
	}

	return result, nil
}

type Stat struct {
	//TODO: Add VisitorsToday
	ParticipantsTotal int
	ParticipantsToday int

	RoundStage       mathbattle.RoundStage
	TimeToSolveLeft  time.Duration
	TimeToReviewLeft time.Duration
	SolutionsTotal   int
	ReviewsTotal     int
	// TODO: Add SolutionsToday
	// TODO: Add ReviewsToday
}

func StatReport(participants mathbattle.ParticipantRepository, rounds mathbattle.RoundRepository,
	solutions mathbattle.SolutionRepository, reviews mathbattle.ReviewRepository) (Stat, error) {

	result := Stat{}

	pAll, err := participants.GetAll()
	if err != nil {
		return result, err
	}
	pToday := mathbattle.FilterRegisteredAfter(pAll, time.Now().Truncate(24*time.Hour))
	result.ParticipantsTotal = len(pAll)
	result.ParticipantsToday = len(pToday)

	round, err := rounds.GetRunning()
	result.RoundStage = mathbattle.GetRoundStage(round)
	if err != nil {
		if err != mathbattle.ErrNotFound {
			return result, err
		}
	} else {
		result.TimeToSolveLeft = time.Until(round.GetSolveEndDate())
		result.TimeToReviewLeft = time.Until(round.GetReviewEndDate())

		allRoundSolutions, err := solutions.FindMany(round.ID, "", "")
		if err != nil {
			return result, err
		}
		result.SolutionsTotal = len(allRoundSolutions)

		allRoundReviews, err := reviews.FindMany(round.ID, "")
		if err != nil {
			return result, err
		}
		result.ReviewsTotal = len(allRoundReviews)
	}

	return result, nil
}
