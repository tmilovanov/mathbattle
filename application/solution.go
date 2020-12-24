package application

import (
	"mathbattle/models/mathbattle"
)

type SolutionService struct {
	Rep    mathbattle.SolutionRepository
	Rounds mathbattle.RoundRepository
}

func (s *SolutionService) Create(solution mathbattle.Solution) (mathbattle.Solution, error) {
	return s.Rep.Store(solution)
}

func (s *SolutionService) Get(ID string) (mathbattle.Solution, error) {
	return s.Rep.Get(ID)
}

func (s *SolutionService) Find(descriptor mathbattle.FindDescriptor) ([]mathbattle.Solution, error) {
	return s.Rep.FindMany(descriptor.RoundID, descriptor.ParticipantID, descriptor.ProblemID)
}

func (s *SolutionService) AppendPart(ID string, part mathbattle.Image) error {
	return s.Rep.AppendPart(ID, part)
}

func (s *SolutionService) Update(solution mathbattle.Solution) error {
	return s.Rep.Update(solution)
}

func (s *SolutionService) Delete(ID string) error {
	return s.Rep.Delete(ID)
}

func (s *SolutionService) GetProblemDescriptors(participantID string) ([]mathbattle.ProblemDescriptor, error) {
	result := []mathbattle.ProblemDescriptor{}

	round, err := s.Rounds.GetRunning()
	if err != nil {
		return result, err
	}

	for _, desc := range round.ProblemDistribution[participantID] {
		result = append(result, desc)
	}

	return result, nil
}
