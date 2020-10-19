package problemdistributor

import (
	"fmt"
	mathbattle "mathbattle/models"
)

// SimpleDistributor gives each participant first problem that is suitable to participant by grade
type SimpleDistributor struct {
	Problems mathbattle.ProblemRepository
}

func (d *SimpleDistributor) GetProblemsForParticipant(participant mathbattle.Participant, count int) ([]mathbattle.Problem, error) {
	result := []mathbattle.Problem{}

	allProblems, err := d.Problems.GetAll()
	if err != nil {
		return result, nil
	}

	for i := 0; i < len(allProblems) && len(result) < count; i++ {
		if mathbattle.IsProblemSuitableForParticipant(&allProblems[i], &participant) {
			result = append(result, allProblems[i])
		}
	}

	if len(result) < count {
		return result, fmt.Errorf("Not enough suitable problems for this participant %v", participant)
	}

	return result, nil
}
