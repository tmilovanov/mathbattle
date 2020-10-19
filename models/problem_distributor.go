package models

type ProblemDistributor interface {
	GetForParticipant(participant Participant) ([]Problem, error)
	GetForParticipantCount(participant Participant, count int) ([]Problem, error)
}

func IsProblemSuitableForParticipant(problem *Problem, participant *Participant) bool {
	if participant.Grade >= problem.MinGrade && participant.Grade <= problem.MaxGrade {
		return true
	}
	return false
}
