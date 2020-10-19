package models

type ProblemDistributor interface {
	GetProblemsForParticipant(participant Participant, count int) ([]Problem, error)
}

func IsProblemSuitableForParticipant(problem *Problem, participant *Participant) bool {
	if participant.Grade >= problem.MinGrade && participant.Grade <= problem.MaxGrade {
		return true
	}
	return false
}
