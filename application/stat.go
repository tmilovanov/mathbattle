package application

import (
	"time"

	"mathbattle/models/mathbattle"
)

type StatService struct {
	Participants mathbattle.ParticipantRepository
	Rounds       mathbattle.RoundRepository
	Solutions    mathbattle.SolutionRepository
	Reviews      mathbattle.ReviewRepository
}

func FilterRegisteredAfter(participants []mathbattle.Participant, datetime time.Time) []mathbattle.Participant {
	result := []mathbattle.Participant{}

	for _, participant := range participants {
		if participant.RegistrationTime.After(datetime) {
			result = append(result, participant)
		}
	}

	return result
}

func (ss *StatService) Stat() (mathbattle.Stat, error) {
	result := mathbattle.Stat{}

	pAll, err := ss.Participants.GetAll()
	if err != nil {
		return result, err
	}
	pToday := FilterRegisteredAfter(pAll, time.Now().Truncate(24*time.Hour))
	result.ParticipantsTotal = len(pAll)
	result.ParticipantsToday = len(pToday)

	round, err := ss.Rounds.GetRunning()
	result.RoundStage = mathbattle.GetRoundStage(round)
	if err != nil {
		if err != mathbattle.ErrNotFound {
			return result, err
		}
	} else {
		result.TimeToSolveLeft = time.Until(round.GetSolveEndDate())
		result.TimeToReviewLeft = time.Until(round.GetReviewEndDate())

		allRoundSolutions, err := ss.Solutions.FindMany(round.ID, "", "")
		if err != nil {
			return result, err
		}
		result.SolutionsTotal = len(allRoundSolutions)

		allRoundReviews, err := ss.Reviews.FindMany(round.ID, "")
		if err != nil {
			return result, err
		}
		result.ReviewsTotal = len(allRoundReviews)
	}

	return result, nil
}
