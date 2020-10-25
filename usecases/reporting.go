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
	for _, solution := range d.ToOrganizers {
		p, err := participants.GetByID(solution.ParticipantID)
		if err != nil {
			return "", err
		}

		result += fmt.Sprintf("'%s', %d grade, solution on %s\n", p.Name, p.Grade, solution.ProblemID)
	}
	result += "\n"

	result += "Between participants: \n"
	result += "---------\n"
	for solutionID, participantIDs := range d.BetweenParticipants {
		solution, err := solutions.Get(solutionID)
		if err != nil {
			return "", err
		}
		p, err := participants.GetByID(solution.ParticipantID)
		if err != nil {
			return "", err
		}
		for _, targetParticipantID := range participantIDs {
			targetParticipant, err := participants.GetByID(targetParticipantID)
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
