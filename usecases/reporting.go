package usecases

import (
	"fmt"
	mathbattle "mathbattle/models"
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
			targetParticipant, err := participants.GetByID(solution.ParticipantID)
			if err != nil {
				return "", err
			}
			result += fmt.Sprintf("'%s'(Problem: %s) -> '%s'\n", p.Name, solution.ProblemID, targetParticipant.Name)
		}
	}

	return result, nil
}

func StatReport(participants mathbattle.Participant, solutions mathbattle.Solution, rounds mathbattle.Round) (string, error) {
	return "", nil
}
