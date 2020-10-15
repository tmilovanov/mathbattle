package models

type ProblemDistributor interface {
	Get(participants []Participant, problems []Problem, rounds []Round) (RoundDistribution, error)
}

func IsProblemSuitableForParticipant(problem *Problem, participant *Participant) bool {
	if participant.Grade >= problem.MinGrade && participant.Grade <= problem.MaxGrade {
		return true
	}
	return false
}
