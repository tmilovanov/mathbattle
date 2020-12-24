package application

import "mathbattle/models/mathbattle"

type ProblemService struct {
	Rep mathbattle.ProblemRepository
}

func (s *ProblemService) GetByID(ID string) (mathbattle.Problem, error) {
	return s.Rep.GetByID(ID)
}
