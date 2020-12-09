package problemdistributor

import (
	"fmt"

	"mathbattle/containers"
	mathbattle "mathbattle/models"
	"mathbattle/mstd"

	"github.com/pkg/errors"
)

type RandomDistributor struct{}

func isProblemAlreadyUsed(participant mathbattle.Participant, problem mathbattle.Problem, pastRounds []mathbattle.Round) bool {
	for _, round := range pastRounds {
		_, err := round.ProblemDistribution.FindDescriptor(participant.ID, problem.ID)
		if err != nil {
			continue
		}

		return true
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

	var result mathbattle.RoundDistribution = make(map[string][]mathbattle.ProblemDescriptor)
	counter := containers.NewUsageCounter()

	for _, participant := range participants {
		suitableProblems := getSuitableProblems(participant, problems, rounds)
		if len(suitableProblems) < count {
			return result, fmt.Errorf("Not enough suitable problems for this participant %v", participant)
		}

		problemIDs := mathbattle.GetProblemIDs(suitableProblems)

		counter.AddItems(problemIDs)

		problemsToSend := counter.UseMostUnpopularFromSet(problemIDs, count)
		for i, problemID := range problemsToSend {
			result[participant.ID] = append(result[participant.ID], mathbattle.ProblemDescriptor{
				Caption:   mstd.IndexToLetter(i),
				ProblemID: problemID,
			})
		}
	}

	return result, nil
}

func (d *RandomDistributor) GetProblemsForParticipant(participantID string, count int) ([]mathbattle.Problem, error) {
	return []mathbattle.Problem{}, errors.Errorf("Not implemented")
}
