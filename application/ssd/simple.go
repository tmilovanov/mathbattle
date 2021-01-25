// ssd означает Solve Stage Distributor, распределитель задач между участниками, на этапе решения задач (Solve Stage)
package ssd

import (
	"fmt"

	"mathbattle/models/mathbattle"
)

// SimpleDistributor gives each participant first problem that is suitable to participant by grade
type SimpleDistributor struct {
	problems     mathbattle.ProblemRepository
	defaultCount int
}

func NewSimpleDistributor(problems mathbattle.ProblemRepository, defaultProblemsCount int) SimpleDistributor {
	return SimpleDistributor{
		problems:     problems,
		defaultCount: defaultProblemsCount,
	}
}

func (d *SimpleDistributor) GetForParticipant(participant mathbattle.Participant) ([]mathbattle.Problem, error) {
	return d.GetForParticipantCount(participant, d.defaultCount)
}

func (d *SimpleDistributor) GetForParticipantCount(participant mathbattle.Participant, count int) ([]mathbattle.Problem, error) {
	result := []mathbattle.Problem{}

	allProblems, err := d.problems.GetAll()
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
