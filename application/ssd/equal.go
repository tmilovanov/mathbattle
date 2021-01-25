// ssd означает Solve Stage Distributor, распределитель задач между участниками, на этапе решения задач (Solve Stage)
package ssd

import (
	"errors"
	"mathbattle/models/mathbattle"
)

//EqualDistributor gives each participant exactly the same problems
type EqualDistributor struct {
	problemsToGive []mathbattle.Problem
}

func NewEqualDistributor(problems mathbattle.ProblemRepository, problemsIDs []string) (*EqualDistributor, error) {
	result := &EqualDistributor{}

	if len(problemsIDs) == 0 {
		return nil, errors.New("Can't create distributor that didn't give any problem")
	}

	for _, problemID := range problemsIDs {
		problem, err := problems.GetByID(problemID)
		if err != nil {
			return nil, err
		}
		result.problemsToGive = append(result.problemsToGive, problem)
	}

	return result, nil
}

func (d *EqualDistributor) GetForParticipant(mathbattle.Participant) ([]mathbattle.Problem, error) {
	return d.problemsToGive, nil
}
