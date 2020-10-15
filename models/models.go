package models

type Storage struct {
	Participants ParticipantRepository
	Problems     ProblemRepository
	Solutions    SolutionRepository
	Rounds       RoundRepository
}
