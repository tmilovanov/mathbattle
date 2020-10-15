package reviewdistributor

import (
	mathbattle "mathbattle/models"
)

type OnReviewDistributor struct{}

func MapSolutionsToParticipants(solutions []mathbattle.Solution, count uint) map[string][]string {
	// Maps solutionID to participants it needs to be sent
	result := make(map[string][]string)
	targets := append(solutions[1:], solutions[:count]...)
	for i := 0; i < len(solutions); i++ {
		for j := uint(0); j < count; j++ {
			result[solutions[i].ID] = append(result[solutions[i].ID], targets[uint(i)+j].ParticipantID)
		}
	}
	return result
}

func (d *OnReviewDistributor) Get(allRoundSolutions []mathbattle.Solution, count uint) mathbattle.ReviewDistribution {
	result := mathbattle.ReviewDistribution{
		BetweenParticipants: make(map[string][]string),
		ToOrganizers:        make([]mathbattle.Solution, 0),
	}
	for _, problemSolutions := range mathbattle.SplitInGroupsByProblem(allRoundSolutions) {
		if len(problemSolutions) == 0 {
			continue
		}

		if len(problemSolutions) == 1 {
			result.ToOrganizers = append(result.ToOrganizers, problemSolutions[0])
			continue
		}

		k := uint(count)
		if uint(len(problemSolutions)) < count+1 {
			k = uint(len(problemSolutions)) - 1
		}

		for solutionID, participantIDs := range MapSolutionsToParticipants(problemSolutions, k) {
			result.BetweenParticipants[solutionID] = participantIDs
		}
	}

	return result
}
