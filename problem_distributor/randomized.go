package problemdistributor

import (
	"fmt"

	"mathbattle/containers"
	mathbattle "mathbattle/models"

	"github.com/pkg/errors"
)

type RandomDistributor struct{}

func isFound(problems []string, problemID string) bool {
	for i := 0; i < len(problems); i++ {
		if problems[i] == problemID {
			return true
		}
	}
	return false
}

func isProblemAlreadyUsed(participant mathbattle.Participant, problem mathbattle.Problem, pastRounds []mathbattle.Round) bool {
	for _, round := range pastRounds {
		problemIDs, isExist := round.ProblemDistribution[participant.ID]
		if !isExist { // User didn't participated in this round
			continue
		}

		if isFound(problemIDs, problem.ID) {
			return true
		}
	}

	return false
}

func getSuitableProblems(participant mathbattle.Participant, problems []mathbattle.Problem,
	pastRounds []mathbattle.Round) []mathbattle.Problem {

	result := []mathbattle.Problem{}
	for _, problem := range problems {
		if !mathbattle.IsProblemSuitableForParticipant(&problem, &participant) ||
			isProblemAlreadyUsed(participant, problem, pastRounds) {
			continue
		}

		result = append(result, problem)
	}
	return result
}

func (d *RandomDistributor) GetForAll(participants []mathbattle.Participant, problems []mathbattle.Problem, rounds []mathbattle.Round, count int) (mathbattle.RoundDistribution, error) {

	var result mathbattle.RoundDistribution = make(map[string][]string)
	counter := containers.NewUsageCounter()

	for _, participant := range participants {
		suitableProblems := getSuitableProblems(participant, problems, rounds)
		if len(suitableProblems) < count {
			return result, fmt.Errorf("Not enough suitable problems for this participant %v", participant)
		}

		problemIDs := mathbattle.GetProblemIDs(suitableProblems)

		counter.AddItems(problemIDs)

		result[participant.ID] = counter.UseMostUnpopularFromSet(problemIDs, count)
	}

	return result, nil
}

func (d *RandomDistributor) GetProblemsForParticipant(participantID string, count int) ([]mathbattle.Problem, error) {
	return []mathbattle.Problem{}, errors.Errorf("Not implemented")
}
