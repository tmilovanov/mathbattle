package solutiondistributor

import (
	"math/rand"
	"mathbattle/models/mathbattle"
	"time"
)

type SolutionDistributor struct{}

func distributeSolutionsToParticipants(solutions []mathbattle.Solution, reviewerCount uint) map[string][]string {
	// Shuffle
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(solutions), func(i, j int) {
		solutions[i], solutions[j] = solutions[j], solutions[i]
	})

	// Map each solution ID to participant IDs it needs to be sent
	result := make(map[string][]string)
	targets := append(solutions[1:], solutions[:reviewerCount]...)

	for i := 0; i < len(solutions); i++ {
		for j := uint(0); j < reviewerCount; j++ {
			result[solutions[i].ID] = append(result[solutions[i].ID], targets[uint(i)+j].ParticipantID)
		}
	}

	return result
}

func (d *SolutionDistributor) Get(allRoundSolutions []mathbattle.Solution, reviewerCount uint) mathbattle.ReviewDistribution {
	result := mathbattle.ReviewDistribution{
		BetweenParticipants: make(map[string][]string),
		ToOrganizers:        make([]string, 0),
	}

	for _, problemSolutions := range mathbattle.SplitInGroupsByProblem(allRoundSolutions) {
		if len(problemSolutions) == 0 {
			continue
		}

		if len(problemSolutions) == 1 {
			result.ToOrganizers = append(result.ToOrganizers, problemSolutions[0].ID)
			continue
		}

		finalReviewerCount := reviewerCount
		if uint(len(problemSolutions)) < reviewerCount+1 {
			finalReviewerCount = uint(len(problemSolutions)) - 1
		}

		for solutionID, participantIDs := range distributeSolutionsToParticipants(problemSolutions, finalReviewerCount) {
			for _, pID := range participantIDs {
				result.BetweenParticipants[pID] = append(result.BetweenParticipants[pID], solutionID)
			}
		}
	}

	return result
}
